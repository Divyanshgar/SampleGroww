package example

import (
	"fmt"
	"log"
	"sort"

	"notification-server/config"
	"notification-server/database"
	"notification-server/models"

	"github.com/xuri/excelize/v2"
)

// StockData represents stock data for Excel export
type StockData struct {
	SNo         int
	StockCode   string
	Description string
	Location    string
	Bin         string
	Capacity    string
	Current     int
	Variance    string
	Notes       string
}

func SimpleExcelExample() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Fetch one user for validated by
	var user models.User
	if err := database.DB.First(&user).Error; err != nil {
		log.Printf("No user found for validated by, using default: %v", err)
		user.FirstName = "John"
		user.LastName = "Smith"
	}
	validatedBy := fmt.Sprintf("Validated by %s %s", user.FirstName, user.LastName)

	// Fetch stocks from database
	var stocks []models.Stock
	if err := database.DB.Find(&stocks).Error; err != nil {
		log.Fatalf("Failed to fetch stocks: %v", err)
	}

	fmt.Printf("Fetched %d stocks from database\n", len(stocks))

	// Fetch users for user details
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	fmt.Printf("Fetched %d users from database\n", len(users))

	// Convert to StockData slice
	stockData := make([]StockData, len(stocks))
	for i, stock := range stocks {
		stockData[i] = StockData{
			SNo:         i + 1,
			StockCode:   fmt.Sprintf("%d", stock.ID),
			Description: stock.StockName,
			Location:    "N/A",
			Bin:         "N/A",
			Capacity:    "N/A",
			Current:     stock.Quantity,
			Variance:    "N/A",
			Notes:       stock.Action,
		}
	}

	// Sort stocks by Description (StockName) ascending
	sort.Slice(stockData, func(i, j int) bool {
		return stockData[i].Description < stockData[j].Description
	})

	// Sort users by FirstName ascending
	sort.Slice(users, func(i, j int) bool {
		return users[i].FirstName < users[j].FirstName
	})

	// Create a new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Rename first sheet to "Report"
	f.SetSheetName("Sheet1", "Report")

	// Add image to A1 in Report sheet
	var err error
	logoPath := "C:/Users/DakshBisht/SampleGroww/services/logo2.png"
	err = f.AddPicture("Report", "A1", logoPath, &excelize.GraphicOptions{
		AltText: "Logo",
		ScaleX:  0.5,
		ScaleY:  0.5,
	})
	if err != nil {
		fmt.Printf("Failed to add image: %v\n", err)
	} else {
		fmt.Println("Image added successfully to Report sheet")
	}

	// Add header rows for stock section
	f.SetCellValue("Report", "A2", "Stock Sheet")
	f.SetCellValue("Report", "A3", "ABC Traders Pvt Ltd")
	f.SetCellValue("Report", "B3", "Warehouse: ABC Traders")
	f.SetCellValue("Report", "A4", validatedBy)
	f.MergeCell("Report", "A2", "I2") // Merge for title
	f.MergeCell("Report", "A3", "I3") // Merge for company

	// Add stock table headers starting from row 5
	stockHeaders := []string{"S.No", "Stock Code", "Description", "Location", "Bin", "Capacity", "Current", "Variance", "Notes"}
	for i, header := range stockHeaders {
		cell := fmt.Sprintf("%c5", 'A'+i)
		f.SetCellValue("Report", cell, header)
	}

	// Add sorted stock data starting from row 6
	for _, stock := range stockData {
		row := stock.SNo + 5
		f.SetCellValue("Report", fmt.Sprintf("A%d", row), stock.SNo)
		f.SetCellValue("Report", fmt.Sprintf("B%d", row), stock.StockCode)
		f.SetCellValue("Report", fmt.Sprintf("C%d", row), stock.Description)
		f.SetCellValue("Report", fmt.Sprintf("D%d", row), stock.Location)
		f.SetCellValue("Report", fmt.Sprintf("E%d", row), stock.Bin)
		f.SetCellValue("Report", fmt.Sprintf("F%d", row), stock.Capacity)
		f.SetCellValue("Report", fmt.Sprintf("G%d", row), stock.Current)
		f.SetCellValue("Report", fmt.Sprintf("H%d", row), stock.Variance)
		f.SetCellValue("Report", fmt.Sprintf("I%d", row), stock.Notes)
	}

	// Calculate starting row for user section (after stock data + 2 rows gap)
	userStartRow := 6 + len(stockData) + 2

	// Add title for user section
	f.SetCellValue("Report", fmt.Sprintf("A%d", userStartRow), "User Profile Details")
	f.MergeCell("Report", fmt.Sprintf("A%d", userStartRow), fmt.Sprintf("F%d", userStartRow))

	// Add user headers starting from next row
	userStartRow++
	userHeaders := []string{"First Name", "Last Name", "Email", "Phone", "City", "Country"}
	for i, header := range userHeaders {
		cell := fmt.Sprintf("%c%d", 'A'+i, userStartRow)
		f.SetCellValue("Report", cell, header)
	}

	// Add sorted user data starting from next row
	userStartRow++
	for i, u := range users {
		row := userStartRow + i
		f.SetCellValue("Report", fmt.Sprintf("A%d", row), u.FirstName)
		f.SetCellValue("Report", fmt.Sprintf("B%d", row), u.LastName)
		f.SetCellValue("Report", fmt.Sprintf("C%d", row), u.Email)
		f.SetCellValue("Report", fmt.Sprintf("D%d", row), u.Phone)
		f.SetCellValue("Report", fmt.Sprintf("E%d", row), u.City)
		f.SetCellValue("Report", fmt.Sprintf("F%d", row), u.Country)
	}

	// Save the file
	if err := f.SaveAs("stock_sheet.xlsx"); err != nil {
		fmt.Printf("Failed to save file: %v\n", err)
		return
	}

	fmt.Println("Excel file 'stock_sheet.xlsx' created successfully with stock and user profile details in the same 'Report' sheet.")
}
