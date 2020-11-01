/*
 *
 * Licensed Material - Property of SRJ Consulting (C) Copyright SRJ Consulting 2020
 * All Rights Reserved.
 *
 */

package main

import (
    "encoding/json"
    "encoding/xml"
    "flag"
    "io/ioutil"
    "os"

    "github.com/srjinatl/pfsense-untangle-static-dhcp/log"
    "go.uber.org/zap"
)

const (
    applicationName = "PfSense to Untangle Static DHCP Converter"
)


var version string

type dhcpEntry struct {
    macAddress string
    address string
    description string
}

func main() {
    // process any command line flags
    devLogger := flag.Bool("devlogger", false, "Specify flag to turn on development level logging")
    inputFile := flag.String("inputFile", "data/pfsense.xml", "Input pFSense XML static DHCP export file")
    outputFile := flag.String("outputFile", "data/untangle.json", "Output Untangle JSON static DHCP import file")
    flag.Parse()

    // set up logger
    logger := log.NewLogger(applicationName, *devLogger)
    logger.Zap.Info("Starting application", zap.String("version", version))
    defer logger.Zap.Info("Ending application")

    logger.Zap.Debug("Files to process", zap.Stringp("InputFile", inputFile), zap.Stringp("OutputFile", outputFile))

    // read in the pfsense dhcp list
    dhcpPfEntries, err := readInputFile(*inputFile)
    if err != nil {
        logger.Zap.Fatal("Error reading in pfsense dhcp export xml file", zap.Error(err))
    }
    logger.Zap.Debug("Loaded pfsense xml input file successfully", zap.Int("Num_Entries", len(dhcpPfEntries)))

    // create untangle list
    dhcpUntangleEntries := generateUntangleList(dhcpPfEntries)

    // write out the untangle list
    err = generateUntangleImportFile(*outputFile, dhcpUntangleEntries)
    if err != nil {
        logger.Zap.Fatal("Error writing out Untangle import file", zap.Error(err))
    }

}

type dhcpd struct {
    XMLName xml.Name `xml:"dhcpd"`
    Lans []lan `xml:"lan"`
}

type lan struct {
    XMLName xml.Name `xml:"lan"`
    Entries []staticMap `xml:"staticmap"`
}

type staticMap struct {
    XMLName xml.Name `xml:"staticmap"`
    Mac string  `xml:"mac"`
    IpAddr string `xml:"ipaddr"`
    HostName string `xml:"hostname"`
    Desc string `xml:"descr"`
}

func readInputFile(fileName string) ([]dhcpEntry, error) {
    // Open our xmlFile
    xmlFile, err := os.Open(fileName)
    // if we os.Open returns an error then handle it
    if err != nil {
        return nil, err
    }

    // defer the closing of our xmlFile so that we can parse it later on
    defer func() {
       _ = xmlFile.Close()
    }()

    // read our opened xmlFile as a byte array.
    byteValue, _ := ioutil.ReadAll(xmlFile)

    // we initialize our Users array
    var dhcpd dhcpd
    // we unmarshal our byteArray which contains our
    // xmlFiles content into 'users' which we defined above
    err = xml.Unmarshal(byteValue, &dhcpd)
    if err != nil {
        return nil, err
    }

    var entryList []dhcpEntry
    if len(dhcpd.Lans) > 0 {
        for _, e := range dhcpd.Lans[0].Entries {
            entryList = append(entryList, dhcpEntry{
                macAddress:  e.Mac,
                address:     e.IpAddr,
                description: e.Desc,
            })
        }
    }

    return entryList, nil
}

type untangleDHCPEntry struct {
    MAC string `json:"macAddress"`
    IPAddr string `json:"address"`
    JavaClass string `json:"javaClass"`
    Desc string `json:"description"`
}

func generateUntangleList(dhcpPfEntries []dhcpEntry) []untangleDHCPEntry {
    var untangleEntries []untangleDHCPEntry
    for _, e := range dhcpPfEntries {
        untangleEntries = append(untangleEntries, untangleDHCPEntry{
            MAC:        e.macAddress,
            IPAddr:     e.address,
            JavaClass:  "com.untangle.uvm.network.DhcpStaticEntry",
            Desc:       e.description,
        })
    }
    return untangleEntries
}

func generateUntangleImportFile(outputFile string, dhcpUntangleEntries []untangleDHCPEntry) (err error){
    file, _ := json.MarshalIndent(dhcpUntangleEntries, "", " ")
    err = ioutil.WriteFile(outputFile, file, 0644)
    return
}