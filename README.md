# Dependencies

This software uses the [Go language](https://golang.org/), the [GNU Scientific Library](http://www.gnu.org/software/gsl/) and [matplotlib](http://matplotlib.sourceforge.net/). To obtain the required packages on a recent version of Linux Mint / Ubuntu, run:

    sudo apt-get install golang libgsl-dev python-matplotlib python-tk

Note - for older versions of Debian-based distributions, instead of installing libgsl-dev, obtain GSL with:

    sudo apt-get install gsl-bin libgsl0ldbl libgsl0-dev

For Fedora the corresponding packages are:

    gsl gsl-devel python-matplotlib

The root-finding implementation in this software uses cgo to interface with GSL.
In Go 1.6, a new rule was implemented which forbids passing a Go object
through C which contains a pointer to another Go object.
For now, we need to turn off this rule.
Add the following to ~/.bashrc:

    export GODEBUG=cgocheck=0

# Usage

Data plots are currently built with test scripts. To run all of these scripts,
run allPlot in the root directory. This will create plots in the subdirectories
(tempZero, tempPair, tempCrit, tempFluc).

# References

[S. K. Sarker, PRB 77, 052505 (2008)](http://prb.aps.org/abstract/PRB/v77/i5/e052505)  
[S. K. Sarker and T. Lovorn, PRB 82, 014504 (2010)](http://prb.aps.org/abstract/PRB/v82/i1/e014504)  
[S. K. Sarker and T. Lovorn, PRB 85, 144502 (2012)](http://prb.aps.org/abstract/PRB/v85/i14/e144502)

WebSocket implementation inspired by [chat-websocket-dart](https://github.com/financeCoding/chat-websocket-dart) and [go-websocket-sample](https://github.com/ukai/go-websocket-sample).
