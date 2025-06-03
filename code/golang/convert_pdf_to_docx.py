import argparse
import os
from pdf2docx import Converter
import sys
import logging

logger = logging.getLogger(__name__)
handler = logging.StreamHandler(sys.stdout)
formatter = logging.Formatter("[%(levelname)s] %(message)s")
handler.setFormatter(formatter)
logger.addHandler(handler)
logger.setLevel(logging.INFO)

def convert_pdf_to_docx(pdf_file, output_dir):
    if not os.path.exists(pdf_file):
        logger.error(f"File '{pdf_file}' không tồn tại.")
        return

    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    base_name = os.path.splitext(os.path.basename(pdf_file))[0]
    output_path = os.path.join(output_dir, base_name + ".docx")

    logger.info(f"Đang chuyển đổi '{pdf_file}' sang '{output_path}'...")

    try:
        cv = Converter(pdf_file)
        cv.convert(output_path, start=0, end=None)
        cv.close()
        logger.info("✅ Chuyển đổi hoàn tất.")
    except Exception as e:
        logger.exception(f"Lỗi trong quá trình chuyển đổi: {e}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="Chuyển đổi PDF sang DOCX.")
    parser.add_argument("-f", "--file", required=True, help="Đường dẫn tới file PDF đầu vào")
    parser.add_argument("-o", "--outdir", required=True, help="Thư mục để lưu file DOCX")

    args = parser.parse_args()
    convert_pdf_to_docx(args.file, args.outdir)

if __name__ == "__main__":
    main()
