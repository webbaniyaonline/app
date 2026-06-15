package handlers

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"template/database"

	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/geoip2-golang"
	"github.com/wcharczuk/go-chart/v2"
)

// Page for create country graph
// Fetch IP addresses from PostgreSQL using Raw SQL
func getIPsFromDB() ([]string, error) {
	var ips []string

	// Perform raw query
	rows, err := database.DB.Db.Raw("SELECT ip FROM transaction GROUP BY ip").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Read rows
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}

	return ips, nil
}

// Get country name from IP
func getCountryFromIP(dbReader *geoip2.Reader, ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}

	record, err := dbReader.Country(parsedIP)
	if err != nil {
		log.Printf("Error fetching country for IP %s: %v", ip, err)
		return ""
	}

	return record.Country.Names["en"] // Return the country name in English
}

// function for generate graph by country
func AdminCountryGraph(c *fiber.Ctx) error {

	// Fetch IP addresses from the database using Raw SQL
	ips, err := getIPsFromDB()
	if err != nil {
		return c.Status(500).SendString("Error fetching IPs from database")
	}

	// Load GeoIP database
	dbReader, err := geoip2.Open("views/uploads/GeoLite2-Country.mmdb")
	if err != nil {
		return c.Status(500).SendString("Error loading GeoLite2 database")
	}
	defer dbReader.Close()

	// Map IPs to countries and count occurrences
	countryCounts := map[string]int{}
	for _, ip := range ips {
		country := getCountryFromIP(dbReader, ip)
		if country != "" {
			countryCounts[country]++
		}
	}

	// Generate graph data
	countries := []string{}
	counts := []float64{}
	countriesx := []string{}
	for country, count := range countryCounts {
		countries = append(countries, country)
		counts = append(counts, float64(count))
		//cntx := string(count)
		cntx := strconv.Itoa(count)
		countriesx = append(countriesx, "{ latLng : [35.8617, 104.1954], name : '"+country+" :  "+cntx+"'},")

	}
	fmt.Println("countries => ", countries)
	fmt.Println("countriesx => ", countriesx)
	// Create bar graph
	graph := chart.BarChart{
		Title: "Transaction by Country",
		//Width:  800,
		//Height: 600,
		Bars: []chart.Value{},
	}

	for i, country := range countries {
		graph.Bars = append(graph.Bars, chart.Value{
			Value: counts[i],
			Label: country,
		})
	}

	// Render chart to a buffer
	buffer := &strings.Builder{}
	if err := graph.Render(chart.PNG, buffer); err != nil {
		return c.Status(500).SendString("Error rendering graph")
	}

	// Send graph as PNG
	c.Set("Content-Type", "image/png")

	return c.SendString(buffer.String())

}
