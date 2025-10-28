package services

import (
	"context"
	"fmt"
	"notification-server/database"
	"notification-server/models"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct{}

// NewExcelService creates a new Excel service instance
func NewExcelService() *ExcelService {
	return &ExcelService{}
}

// GenerateUserProfileExcel generates an Excel file with user profile and stock details in a single template sheet
func (s *ExcelService) GenerateUserProfileExcel(user *models.User) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Create single "Profile Template" sheet
	f.SetSheetName("Sheet1", "Profile Template")
	f.SetActiveSheet(0)

	// Add company info to top
	if err := s.addLogo(f); err != nil {
		return nil, fmt.Errorf("failed to add company info: %w", err)
	}

	// Add user details starting from row 4
	if err := s.addUserDetails(f, user); err != nil {
		return nil, fmt.Errorf("failed to add user details: %w", err)
	}

	// Add stock details below user details starting from row 16
	if err := s.addStockDetails(f, user); err != nil {
		return nil, fmt.Errorf("failed to add stock details: %w", err)
	}

	// Set column widths for better visibility
	columns := []string{"A", "B", "C", "D", "E", "F"}
	for _, col := range columns {
		f.SetColWidth("Profile Template", col, col, 25) // Increased width for better visibility
	}

	// Set row heights for better spacing
	f.SetRowHeight("Profile Template", 1, 25)  // Company name row
	f.SetRowHeight("Profile Template", 2, 20)  // Date row
	f.SetRowHeight("Profile Template", 3, 15)  // Spacer
	f.SetRowHeight("Profile Template", 4, 25)  // User title row
	f.SetRowHeight("Profile Template", 5, 15)  // Spacer
	f.SetRowHeight("Profile Template", 16, 25) // Stock title row
	f.SetRowHeight("Profile Template", 17, 15) // Spacer

	// Freeze top row for better navigation
	f.SetPanes("Profile Template", &excelize.Panes{
		Freeze: true,
		YSplit: 1,
	})

	// Arial font will be applied via individual styles in addUserDetails and addStockDetails

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ✅ Updated for latest Excelize version
func (s *ExcelService) addLogo(f *excelize.File) error {
	// Try local logos first, fallback to OneDrive
	logoPaths := []string{
		"./templates/logo.jpg",
		"./services/logo2.png",
		`c:/Users/DakshBisht/OneDrive - CENTRICITY FINANCIAL DISTRIBUTION PRIVATE LIMITED/Pictures/Camera Roll/centricity logo.png`,
	}

	var absPath string
	var err error
	for _, logoPath := range logoPaths {
		if s.fileExists(logoPath) {
			absPath, err = filepath.Abs(logoPath)
			if err != nil {
				continue // Try next path
			}
			break // Found a valid path
		}
	}

	sheet := "Profile Template"

	// Attempt to add the logo image to A1 in "Profile Template" sheet
	if absPath != "" {
		if err := f.AddPicture(sheet, "A1", absPath, &excelize.GraphicOptions{
			AltText: "Centricity Logo",
			ScaleX:  1.0,
			ScaleY:  1.0,
		}); err != nil {
			// Log the error but continue without the logo
			fmt.Printf("Warning: Failed to add logo from %s: %v. Continuing without logo.\n", absPath, err)
		}
	} else {
		fmt.Println("Warning: No valid logo file found. Continuing without logo.")
	}

	// Always add company name in A3 (merged A3:B3), regardless of logo success
	if err := f.MergeCell(sheet, "A3", "B3"); err != nil {
		return fmt.Errorf("failed to merge cells for company name: %w", err)
	}
	f.SetCellValue(sheet, "A3", "Centricity Financial Distribution Private Limited")

	// Company name style
	companyStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Bold:   true,
			Size:   14,
			Color:  "#000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}
	f.SetCellStyle(sheet, "A3", "B3", companyStyle)

	// Always add generation date in A4 (merged A4:B4), regardless of logo success
	if err := f.MergeCell(sheet, "A4", "B4"); err != nil {
		return fmt.Errorf("failed to merge cells for date: %w", err)
	}
	f.SetCellValue(sheet, "A4", fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))

	// Date style
	dateStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   11,
			Color:  "#666666",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}
	f.SetCellStyle(sheet, "A4", "B4", dateStyle)

	return nil
}

// fileExists checks if a file exists and is not a directory
func (s *ExcelService) fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (s *ExcelService) addUserDetails(f *excelize.File, user *models.User) error {
	sheet := "Profile Template"

	// Add section title at A5 (merged A5:B5)
	if err := f.MergeCell(sheet, "A5", "B5"); err != nil {
		return fmt.Errorf("failed to merge cells for user title: %w", err)
	}
	f.SetCellValue(sheet, "A5", "User Profile Summary")

	// Title style
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Bold:   true,
			Size:   16,
			Color:  "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#667eea"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}
	f.SetCellStyle(sheet, "A5", "B5", titleStyle)

	// Create styles for table
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Bold:   true,
			Size:   11,
			Color:  "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#667eea"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}

	// Alternating row styles
	evenRowStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   11,
			Color:  "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F2F2F2"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "middle",
			WrapText:   true,
		},
	})
	if err != nil {
		return err
	}

	oddRowStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Bold:   true,
			Size:   11,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "middle",
			WrapText:   true,
		},
	})
	if err != nil {
		return err
	}

	// User details as table starting from row 7
	userData := [][]string{
		{"Field", "Value"},
		{"User ID", fmt.Sprintf("%d", user.ID)},
		{"First Name", user.FirstName},
		{"Last Name", user.LastName},
		{"Email", user.Email},
		{"Phone", user.Phone},
		{"Address", user.Address},
		{"City", user.City},
		{"Country", user.Country},
	}

	for i, row := range userData {
		rowNum := i + 7 // Start from row 7 after title and spacer
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), row[0])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), row[1])

		if i == 0 {
			f.SetCellStyle(sheet, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("B%d", rowNum), headerStyle)
		} else {
			style := evenRowStyle
			if i%2 == 1 {
				style = oddRowStyle
			}
			f.SetCellStyle(sheet, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("B%d", rowNum), style)
		}
	}

	return nil
}

func (s *ExcelService) addStockDetails(f *excelize.File, user *models.User) error {
	sheet := "Profile Template"

	// Fetch real stock data from database
	dbStocks, err := database.GetQueries().ListStocksByUser(context.Background(), int32(user.ID))
	if err != nil {
		return fmt.Errorf("failed to fetch stocks: %w", err)
	}

	// Convert to models.Stock
	stocks := make([]models.Stock, len(dbStocks))
	for i, dbStock := range dbStocks {
		stocks[i] = database.ConvertDBStockToModel(dbStock)
	}

	// Sort stocks by date
	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].Date.Before(stocks[j].Date)
	})

	// Add section title at A17 (merged A17:F17)
	if err := f.MergeCell(sheet, "A17", "F17"); err != nil {
		return fmt.Errorf("failed to merge cells for stock title: %w", err)
	}
	f.SetCellValue(sheet, "A17", "Stock Transactions")

	// Title style
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Bold:   true,
			Size:   16,
			Color:  "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#667eea"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}
	f.SetCellStyle(sheet, "A17", "F17", titleStyle)

	if len(stocks) == 0 {
		// Add no data message
		if err := f.MergeCell(sheet, "A20", "F20"); err != nil {
			return fmt.Errorf("failed to merge cells for no data message: %w", err)
		}
		f.SetCellValue(sheet, "A20", "No stock transactions available for this user.")
		noDataStyle, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Family: "Arial",
				Size:   11,
				Color:  "#666666",
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "middle",
			},
		})
		if err != nil {
			return err
		}
		f.SetCellStyle(sheet, "A20", "F20", noDataStyle)
		return nil
	}

	// Create styles for table
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Bold:   true,
			Size:   11,
			Color:  "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#667eea"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}

	// Alternating row styles
	evenRowStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   11,
			Color:  "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F2F2F2"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}

	oddRowStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   11,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}

	// Buy/Sell conditional styles
	buyStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   11,
			Color:  "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D4EDDA"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}

	sellStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   11,
			Color:  "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F8D7DA"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}

	// Totals style
	totalsStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Bold:   true,
			Size:   11,
			Color:  "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E9ECEF"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "middle",
		},
	})
	if err != nil {
		return err
	}

	// Header for stocks table starting from row 20
	headers := []string{"Stock Name", "Action", "Quantity", "Price", "Original Value", "Date"}
	startRow := 20

	for i, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+i, startRow)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	var totalValue float64
	for i, stock := range stocks {
		row := startRow + i + 1
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), stock.StockName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stock.Action)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), stock.Quantity)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), stock.Price)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), stock.OriginalValue)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), stock.Date.Format("2006-01-02"))
		totalValue += stock.OriginalValue

		// Apply conditional and alternating styles
		actionStyle := evenRowStyle
		if stock.Action == "buy" {
			actionStyle = buyStyle
		} else if stock.Action == "sell" {
			actionStyle = sellStyle
		}
		if i%2 == 1 {
			actionStyle = oddRowStyle // Override alternating if conditional
		}
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row), actionStyle)
	}

	// Add totals row
	totalsRow := startRow + len(stocks) + 1
	f.SetCellValue(sheet, fmt.Sprintf("D%d", totalsRow), "Total Original Value")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", totalsRow), totalValue)
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", totalsRow), fmt.Sprintf("F%d", totalsRow), totalsStyle)

	return nil
}
