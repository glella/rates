# Simple Utility to download exchange rates from openechangerates.org and save them to CSV file

Downloads rates for the last business day of the month, for each month year to date and previous year.
Saves CSV file in the same directory as executable.

Versions in Python and Go. Soon in Rust.

Comments on Python version:
- Uses pandas library. It needs to be installed first. Can be done with pip like this:

	pip install --upgrade pip
	
	pip install pandas


- Can be run as an interpreted script or as an executable created with pyinstaller.
That is why the code checks for if it is being run frozen (as part of executable) or not to properly check where to save CSV.

- Pyinstaller can be installed with:

	pip install pyinstaller


- Executable can be created easily with:

	pyinstaller rates.py --name rates --onefile


- If running frozen (as part of executable) needs ssl & certifi to be able to access internet due to permissions in other computers. Hope this is the right way to solve this. Not 100% sure.


Comments on Go version:
- It was refreshing to get an executable and it to work flawlessly in any computer with no fuss unlike Python (when the recipient is not a developer and does not have everything setup as we do).

- I am learning Go and this exercise was done to hack a bit in it, so it is not great coding.

- Strugled in particular to get rates from the struct into an array looping through the list of rates I wanted. Resorted to a brute force repeat of the same command +20 times.

- Another bad tecnique is that I invoked time.Now() everywhere instead of passing it around. That was mostly lazyness.

- Other than that it seems to work very well. Instead of relying on libraries like I did in the Python version with pandas, in Go we calculated by hand the dates of the monthly last business days and the transposing of the data (columns for rows).

