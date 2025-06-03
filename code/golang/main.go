package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

func main() {
	os.MkdirAll("output", os.ModePerm)
	os.MkdirAll("converted", os.ModePerm)

	r := gin.Default()

	r.POST("/convert", func(c *gin.Context) {
		url := c.PostForm("url")
		fileType := c.DefaultPostForm("type", "pdf") // default là pdf
		typeExecute := c.DefaultPostForm("type_exec", "python")

		if url == "" {
			c.JSON(400, gin.H{"error": "Missing url parameter"})
			return
		}

		pdfPath := "output/output.pdf"
		docxPath := "converted/output.docx"

		//err := generatePDF(url, pdfPath)
		//if err != nil {
		//	c.JSON(500, gin.H{"error": "Failed to generate PDF", "details": err.Error()})
		//	return
		//}

		if fileType == "doc" {
			err := convertPDFToDocx(pdfPath, docxPath, typeExecute)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to convert to DOCX", "details": err.Error()})
				return
			}
			c.JSON(200, gin.H{
				"type": "docx",
				"docx": docxPath,
				"pdf":  pdfPath, // optional: include original
			})
		} else {
			c.JSON(200, gin.H{
				"type": "pdf",
				"pdf":  pdfPath,
			})
		}
	})

	r.Run(":8080")
}

func generatePDF(url, outputPath string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),                    // chạy chế độ headless
		chromedp.Flag("headless", "new"),                   // dùng headless mới (Chromium 109+)
		chromedp.Flag("no-sandbox", true),                  // bỏ sandbox (bắt buộc trong Docker)
		chromedp.Flag("disable-gpu", true),                 // không cần GPU
		chromedp.Flag("disable-dev-shm-usage", true),       // tránh lỗi bộ nhớ chia sẻ trong Docker
		chromedp.Flag("disable-software-rasterizer", true), // cần thiết nếu không có GPU
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Tạo context trình duyệt
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var buf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(10 * time.Second), // Đợi trang tải xong (bao gồm ảnh)
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0.5). // inch (~12.7mm)
				WithMarginBottom(0.5).
				WithMarginLeft(0.5).
				WithMarginRight(0.5).
				Do(ctx)
			return err
		}),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf, 0644)
}

func convertPDFToDocx(inputPDF, outputDocx, execType string) error {
	var cmd *exec.Cmd
	//currentPath, err := os.Getwd()

	//if err != nil {
	//	return err
	//}
	switch execType {
	case "python":
		//execPath := path.Join(currentPath, "assets", "convert_pdf_to_docx.py")
		cmd = exec.Command("python", "/app/assets/convert_pdf_to_docx.py", "--file", inputPDF, "--outdir", filepath.Dir(outputDocx))
	case "libreoffice":
		cmd = exec.Command("libreoffice", "--headless", "--infilter=writer_pdf_import", "--convert-to", "doc", "--outdir", filepath.Dir(outputDocx), inputPDF)
	case "python_binary":
		cmd = exec.Command("./assets/convert_pdf_to_docx", "--file", inputPDF, "--outdir", filepath.Dir(outputDocx))
	default:
		return fmt.Errorf("unknown exec type: %s", execType)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
