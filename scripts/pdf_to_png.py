import sys
import os
import fitz  # PyMuPDF
from PIL import Image

def pdf_to_png(pdf_path, output_path):
    try:
        # Normalize path for Windows
        pdf_path = os.path.normpath(pdf_path)
        output_path = os.path.normpath(output_path)
        
        # Check if PDF exists
        if not os.path.exists(pdf_path):
            print(f"ERROR: PDF file not found: {pdf_path}", file=sys.stderr)
            return 1
        
        # Open PDF
        doc = fitz.open(pdf_path)
        
        if len(doc) == 0:
            print("ERROR: PDF has no pages", file=sys.stderr)
            doc.close()
            return 1
        
        # Get first page
        page = doc[0]
        
        # Render page to image at 300 DPI
        mat = fitz.Matrix(300/72, 300/72)  # 300 DPI
        pix = page.get_pixmap(matrix=mat)
        
        # Convert to PIL Image
        img = Image.frombytes("RGB", [pix.width, pix.height], pix.samples)
        
        # Save as PNG
        img.save(output_path, 'PNG')
        
        doc.close()
        
        # Print output path for confirmation
        print(output_path)
        return 0
        
    except Exception as e:
        print(f"ERROR: {str(e)}", file=sys.stderr)
        import traceback
        traceback.print_exc(file=sys.stderr)
        return 1

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python pdf_to_png.py <input.pdf> <output.png>", file=sys.stderr)
        sys.exit(1)
    
    pdf_path = sys.argv[1]
    output_path = sys.argv[2]
    
    sys.exit(pdf_to_png(pdf_path, output_path))
