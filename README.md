# pfsense-untangle-static-dhcp
Translator from pfSense to Untangle static DHCP format

This is a simple utility written in Go to convert the 
static DHCP entries that you might have defined in your
pfSense configuration to the format that can be imported
by Untangle.

Note that you may need to adjust the json tags in the *dhcpd*
struct if your network has different names for your internal
networks. Also note that this program currently supports only
pulling one DHCP server entries from the pfSense input
file. If you have multiple DHCP servers for different 
networks then you will need to adapt this code to either
handle multiple networks or manually adjust the json tag
for each network that you want to export in separate 
runs of this utility.

Command format: 

To obtain the information on how to run this utility, 
enter the following: 

On a Mac:

     convert -h

On a Windows machine:

     convert.exe -h

The output from the above should look similar to that
shown below which explains how to invoke this utility:

    Usage of convert:

    -devlogger
     Specify flag to turn on development level logging
    -inputFile string
     Input pFSense XML static DHCP export file (default "data/pfsense.xml")
    -outputFile string
     Output Untangle JSON static DHCP import file (default "data/untangle.json")

Note that the output file will be overwritten if a file
by the name specified already exists.
