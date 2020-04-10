# Simple Utility to download exchange rates from openexchangerates.org and save them to CSV file

Downloads rates for the last business day of the month, for each month year to date and previous year.
Saves them to a CSV file in the same directory as executable.

Versions in Python, Go and Rust.

Wrote them in that exact order:
- Very quickly got it done in Python running as script. Pulled quite a few hairs to get it running smoothly as a frozen executable in somebody else's computer.
- Refreshingly simple to code it in Go despite being new to the language. Took the same time to get it done as in Python.
- Straightforward implementation in Rust although had to pay deeper attention to ownership details. Took twice as long to code compared to Go and Python.

Executable sizes:

Python 30.4 MB (pyinstaller frozen executable)

Go	7.3 MB

Rust	3.6 MB


Comments on Python version:
- Uses pandas library. It needs to be installed first. Can be done with pip like this:

	pip install --upgrade pip
	
	pip install pandas


- Can be run as an interpreted script or as an executable created with pyinstaller.
That is why the code checks for if it is being run frozen (as part of executable) or not, to properly check where to save CSV.

- Pyinstaller can be installed with:

	pip install pyinstaller


- Executable can be created easily with:

	pyinstaller rates.py --name rates --onefile


- If running frozen (as part of executable) needs ssl & certifi to be able to access internet due to permissions in other computers. Hope this is the right way to solve this. Not 100% sure.


Comments on Go version:
- It was refreshing to get an executable and it to work flawlessly in any computer with no fuss unlike Python (when the recipient is not a developer and does not have everything setup as we do).

- Instead of relying on libraries like I did in the Python version with pandas, in Go we calculated by hand the dates of the monthly last business days, and the transposing of the data (columns for rows).

- Language of choice to quickly hack something together because of simplicity.


Comments on Rust version:
- Same approach as in Go. The feeling of "if it compiles, it works - no surprises" is awesome.

- It took me twice as long to code as the Go version despite knowing exactly what I wanted to do.
Example: the transpose fucntion was originally written with &str instead of Strings, but after downloading the data got into ownership/borrowing issues, so I had to solve that by migrating a few functions to use Strings instead of &str in a domino effect.

- Definitely my favorite language but it takes some time to iron everything out. It is very efficient in terms of size and speed. The feeling of its executable being bulletproof is great.



