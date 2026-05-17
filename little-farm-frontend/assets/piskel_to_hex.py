# Replace this with your array from the .c file
piskel_data = []


def to_hex(color):
	# Extracts RR, GG, BB and formats as #RRGGBB
	r = color & 0xFF
	g = (color >> 8) & 0xFF
	b = (color >> 16) & 0xFF
	return f"#{r:02x}{g:02x}{b:02x}"


for color in piskel_data:
	print(f"{hex(color)}  ->  {to_hex(color)}")
