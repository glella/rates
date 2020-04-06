#!/usr/bin/env python3

# Install instructions:
#	pip install --upgrade pip
# 	pip install pandas


import datetime
from urllib.request import urlopen
import csv
import pandas as pd
import os
import sys
import certifi
import ssl


def last_business_days(year, month):
	start = datetime.date(year, 1, 1)
	end = datetime.date(year, month-1, 31)
	dtindex = pd.date_range(start, end, freq='BM').to_pydatetime().tolist()	
	dates = []
	for val in dtindex:
		dates.append(val.date().strftime('%Y-%m-%d'))
	return dates


def download_data(aDate):
	url = 'https://openexchangerates.org/api/historical/' + aDate + '.json'
	# You need your own key for an app_id. Get it at openexchangerates.org
	# Without a key it wont work.
	url += '?app_id=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
	json = urlopen(url, context=ssl.create_default_context(cafile=certifi.where())).read()
	(true,false,null) = (True,False,None)
	data = eval(json)
	return data.get('rates')


def transpose(filename):
	pd.read_csv(filename, header=None).T.to_csv(filename, header=False, index=False)


def resource_path(relative):
	application_path = ""
	if getattr(sys, 'frozen', False):
		# we are running in a bundle
		application_path = os.path.dirname(sys.argv[0])		
	else:
		# we are running in a normal Python environment
		application_path = os.path.dirname(os.path.abspath(__file__))
	
	return os.path.join(application_path, relative)


if __name__ == "__main__":
	today = datetime.date.today()
	current_year = last_business_days(today.year, today.month)
	previous_year = last_business_days(today.year-1, today.month)
	csv_name = 'ytd_rates_downloaded_' + str(today) + '.csv'
	filename = resource_path(csv_name)

	with open(filename,'w') as f:
		w = csv.writer(f)
		w.writerow(['Rates','AED', 'ARS', 'AUD', 'BRL', 'CAD', 'CLP', 'CNY', 'COP', 'EUR', 'GBP', 'IDR', 'INR', 'IQD', 'JOD', 'JPY', 'KES', 'KRW', 'LKR', 'MXN', 'MYR', 'NGN', 'NZD', 'RUB', 'SAR', 'SGD', 'ZAR'])
		for aDate in current_year:
			rates = download_data(aDate)
			w.writerow([str(aDate), rates.get('AED'), rates.get('ARS'), rates.get('AUD'), rates.get('BRL'), rates.get('CAD'), rates.get('CLP'), rates.get('CNY'), rates.get('COP'), rates.get('EUR'), rates.get('GBP'), rates.get('IDR'), rates.get('INR'), rates.get('IQD'), rates.get('JOD'), rates.get('JPY'), rates.get('KES'), rates.get('KRW'), rates.get('LKR'), rates.get('MXN'), rates.get('MYR'), rates.get('NGN'), rates.get('NZD'), rates.get('RUB'), rates.get('SAR'), rates.get('SGD'), rates.get('ZAR')])
		for aDate in previous_year:
			rates = download_data(aDate)
			w.writerow([str(aDate), rates.get('AED'), rates.get('ARS'), rates.get('AUD'), rates.get('BRL'), rates.get('CAD'), rates.get('CLP'), rates.get('CNY'), rates.get('COP'), rates.get('EUR'), rates.get('GBP'), rates.get('IDR'), rates.get('INR'), rates.get('IQD'), rates.get('JOD'), rates.get('JPY'), rates.get('KES'), rates.get('KRW'), rates.get('LKR'), rates.get('MXN'), rates.get('MYR'), rates.get('NGN'), rates.get('NZD'), rates.get('RUB'), rates.get('SAR'), rates.get('SGD'), rates.get('ZAR')])
		f.close()

	transpose(filename)

