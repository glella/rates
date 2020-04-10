// rates.rs
// Downloads rates from openexchangerates.org and saves them to CVS file
// Downloads YTD and PY last business days

extern crate csv; 
use std::error::Error;
use std::env;
use std::process;
extern crate chrono;
use chrono::{DateTime, Local, NaiveDate, Utc, Datelike, Duration};
extern crate serde_json;
extern crate serde;
#[macro_use]
extern crate serde_derive;
extern crate reqwest;

// Macro to transform a vec of literals (&str) into a vec of Strings
macro_rules! vec_of_strings {
    ($($x:expr),*) => (vec![$($x.to_string()),*]);
}

fn main() {
    
    // Determine dates to download - Calc elapsed months & last business days
    let dates = dates_to_download();

    // Download data
    let mut data: Vec<Vec<String>> = Vec::new();
    
    // Set headers for CVS table
    let header = vec_of_strings!["Rates", "AED", "ARS", "AUD", "BRL", "CAD", "CLP", "CNY", "COP", "EUR", "GBP", "IDR", "INR", "IQD", "JOD", "JPY", "KES", "KRW", "LKR", "MXN", "MYR", "NGN", "NZD", "RUB", "SAR", "SGD", "ZAR"];
    data.push(header);

    // Loop for each date
    for date in dates {
        let rec = download_data(&date);
        let rec = match rec {
            Ok(info) => info,
            Err(err) => {
                println!("{}", err);
                process::exit(1);
            },
        };
        // Select rates to save & set presicion of floats to display as strings
        let record = generate_record(&date, rec);
        // Save the record to Vector
        data.push(record);
    }

    // Transpose columns and rows
    let transposed = transpose(data);

    // Save to CSV file & we are done
    if let Err(err) = save_data(transposed) {
        println!("{}", err);
        process::exit(1);
    }
}

fn elapsed_months() -> i32 {
    let now = Utc::now();
    let (_is_common_era, yr) = now.year_ce();
    let year = yr as i32;
    let mut since = NaiveDate::from_ymd(year, 1, 1);
    let now_naive = NaiveDate::from_ymd(year, now.month(), now.day());
    
    let mut months = 0;
    let mut month = since.month();
    while since < now_naive {
        since = since + Duration::days(1);
        let next_month = since.month();
        if next_month != month {
            months += 1;
        }
        month = next_month;
    }
    months
}

fn num_days_month(year: i32, m: u32) -> i64 {
    let end_date: NaiveDate;
    if m == 12 {
        end_date = NaiveDate::from_ymd(year + 1, 1, 1);
    } else {
        end_date = NaiveDate::from_ymd(year, m + 1, 1);
    }
    let num = end_date.signed_duration_since(NaiveDate::from_ymd(year, m, 1))
    .num_days();

    num
}

fn calc_last_business_day(year: i32, month: u32) -> String {
    let last_day = num_days_month(year, month) as u32;
    let mut date = NaiveDate::from_ymd(year, month, last_day);
    let wd = date.weekday();
    if wd.number_from_monday() == 7 {
        date = NaiveDate::from_ymd(year, month, last_day - 2);
    }
    if wd.number_from_monday() == 6 {
        date = NaiveDate::from_ymd(year, month, last_day - 1);
    }

    return date.format("%Y-%m-%d").to_string()
}

fn dates_to_download() -> Vec<String> {
    let elapsed = elapsed_months() as u32;
    let now: DateTime<Local> = Local::now();
    let current_yr = now.year();
    let previous_yr = current_yr - 1;

    let mut result: Vec<String> = Vec::new();
    for i in 1..= elapsed {
        let date = calc_last_business_day(current_yr, i);
        result.push(date);
    }
    for i in 1..= elapsed {
        let date = calc_last_business_day(previous_yr, i);
        result.push(date);
    }

    result
}

fn download_data(date: &str) -> Result<Historical, reqwest::Error> {
    // Set the url
    // You need to replace xxxxx with your own app_id from openexchangerates.org
    let req_url = format!("https://openexchangerates.org/api/historical/{}.json?app_id=xxxxxxxxxxxxxxxxxxxxxxxxx", date);
    // Get the JSON
    let response = reqwest::blocking::get(&req_url)?;
    // Parse JSON as struct using serde
    let record: Historical = response.json()?;
    // return record if no errors
    Ok(record)
}

fn generate_record(date: &str, r: Historical) -> Vec<String> {
    let mut vec: Vec<String> = Vec::new();
    vec.push(date.to_string());
    vec.push(format!("{:.4}", r.rates.aed));
    vec.push(format!("{:.4}", r.rates.ars));
    vec.push(format!("{:.4}", r.rates.aud));
    vec.push(format!("{:.4}", r.rates.brl));
    vec.push(format!("{:.4}", r.rates.cad));
    vec.push(format!("{:.4}", r.rates.clp));
    vec.push(format!("{:.4}", r.rates.cny));
    vec.push(format!("{:.4}", r.rates.cop));
    vec.push(format!("{:.4}", r.rates.eur));
    vec.push(format!("{:.4}", r.rates.gbp));
    vec.push(format!("{:.4}", r.rates.idr));
    vec.push(format!("{:.4}", r.rates.inr));
    vec.push(format!("{:.4}", r.rates.iqd));
    vec.push(format!("{:.4}", r.rates.jod));
    vec.push(format!("{:.4}", r.rates.jpy));
    vec.push(format!("{:.4}", r.rates.kes));
    vec.push(format!("{:.4}", r.rates.krw));
    vec.push(format!("{:.4}", r.rates.lkr));
    vec.push(format!("{:.4}", r.rates.mxn));
    vec.push(format!("{:.4}", r.rates.myr));
    vec.push(format!("{:.4}", r.rates.ngn));
    vec.push(format!("{:.4}", r.rates.nzd));
    vec.push(format!("{:.4}", r.rates.rub));
    vec.push(format!("{:.4}", r.rates.sar));
    vec.push(format!("{:.4}", r.rates.sgd));
    vec.push(format!("{:.4}", r.rates.zar));
    
    vec
}

fn transpose(data: Vec<Vec<String>>) -> Vec<Vec<String>> {
	let xl = data[0].len(); // columns
	let yl = data.len();    // rows
	let mut result: Vec<Vec<String>> = Vec::new();
	for _i in 0..xl {
        result.push(vec!["".to_string(); yl]);
	}
    for i in 0..xl {
		for j in 0..yl {
			result[i][j] = data[j][i].clone(); // Had to copy data because of Strings
		}
	}

	result
}

fn save_data(data: Vec<Vec<String>>) -> Result<(), Box<dyn Error>> {
    // Get path to executable & determine filename
    let now: DateTime<Local> = Local::now();
    let today: String = now.format("%Y-%m-%d").to_string();
    let name = "ytd_rates_downloaded_".to_owned() + &today + ".csv";
    let path = env::current_dir()?;
    let filename = path.join(name);

    // Create writer
    let mut wtr = csv::Writer::from_path(filename)?;

    // Write data
    for record in data {
        wtr.write_record(record)?;
    }

    wtr.flush()?; // Need to flush buffer
    Ok(())
} 

#[derive(Serialize, Deserialize)]
struct Historical {
    disclaimer: String,
    license: String,
    timestamp: u64,
    base: String,
    rates: Rates,
}

#[derive(Serialize, Deserialize)]
#[serde(rename_all = "UPPERCASE")]
struct Rates {
    aed: f64,
    afn: f64,
    all: f64,
    amd: f64,
    ang: f64,
    aoa: f64,
    ars: f64,
    aud: f64,
    awg: f64,
    azn: f64,
    bam: f64,
    bbd: f64,
    bdt: f64,
    bgn: f64,
    bhd: f64,
    bif: f64,
    bmd: f64,
    bnd: f64,
    bob: f64,
    brl: f64,
    bsd: f64,
    btc: f64,
    btn: f64,
    bwp: f64,
    byn: f64,
    bzd: f64,
    cad: f64,
    cdf: f64,
    chf: f64,
    clf: f64,
    clp: f64,
    cnh: f64,
    cny: f64,
    cop: f64,
    crc: f64,
    cuc: f64,
    cup: f64,
    cve: f64,
    czk: f64,
    djf: f64,
    dkk: f64,
    dop: f64,
    dzd: f64,
    egp: f64,
    ern: f64,
    etb: f64,
    eur: f64,
    fjd: f64,
    fkp: f64,
    gbp: f64,
    gel: f64,
    ggp: f64,
    ghs: f64,
    gip: f64,
    gmd: f64,
    gnf: f64,
    gtq: f64,
    gyd: f64,
    hkd: f64,
    hnl: f64,
    hrk: f64,
    htg: f64,
    huf: f64,
    idr: f64,
    ils: f64,
    imp: f64,
    inr: f64,
    iqd: f64,
    irr: f64,
    isk: f64,
    jep: f64,
    jmd: f64,
    jod: f64,
    jpy: f64,
    kes: f64,
    kgs: f64,
    khr: f64,
    kmf: f64,
    kpw: f64,
    krw: f64,
    kwd: f64,
    kyd: f64,
    kzt: f64,
    lak: f64,
    lbp: f64,
    lkr: f64,
    lrd: f64,
    lsl: f64,
    lyd: f64,
    mad: f64,
    mdl: f64, 
    mga: f64,
    mkd: f64,
    mmk: f64,
    mnt: f64,
    mop: f64,
    mro: f64,
    mru: f64,
    mur: f64,
    mvr: f64,
    mwk: f64,
    mxn: f64,
    myr: f64,
    mzn: f64,
    nad: f64,
    ngn: f64,
    nio: f64,
    nok: f64,
    npr: f64,
    nzd: f64,
    omr: f64,
    pab: f64,
    pen: f64,
    pgk: f64,
    php: f64,
    pkr: f64,
    pln: f64,
    pyg: f64,
    qar: f64,
    ron: f64,
    rsd: f64,
    rub: f64,
    rwf: f64,
    sar: f64,
    sbd: f64,
    scr: f64,
    sdg: f64,
    sek: f64,
    sgd: f64,
    shp: f64,
    sll: f64,
    sos: f64,
    srd: f64,
    ssp: f64,
    std: f64,
    stn: f64,
    svc: f64,
    syp: f64,
    szl: f64,
    thb: f64,
    tjs: f64,
    tmt: f64,
    tnd: f64,
    top: f64,
    #[serde(rename="TRY")]
    tri: f64,
    ttd: f64,
    twd: f64,
    tzs: f64,
    uah: f64,
    ugx: f64,
    usd: f64,
    uyu: f64,
    uzs: f64,
    vef: f64,
    ves: f64,
    vnd: f64,
    vuv: f64,
    wst: f64,
    xaf: f64,
    xag: f64,
    xau: f64,
    xcd: f64,
    xdr: f64,
    xof: f64,
    xpd: f64,
    xpf: f64,
    xpt: f64,
    yer: f64,
    zar: f64,
    zmw: f64,
    zwl: f64,
}

#[cfg(test)]
mod tests {
    use super::*;
 
    #[test]
    fn elapsed_months_works() {
        // in April will suceed - later on should fail
        assert_eq!(elapsed_months(), 3);
    }

    #[test]
    fn num_days_month_works() {
        assert_eq!(num_days_month(2019, 2), 28);
    }

    #[test]
    fn last_bus_day_works() {
        assert_eq!(calc_last_business_day(2019, 3), "2019-03-29");
    }
}
