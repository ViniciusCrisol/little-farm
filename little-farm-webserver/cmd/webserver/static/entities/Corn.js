import { newCellId } from "../utils.js";

export default class Corn {
	id = "";
	kind = "";
	xPosition = -1;
	yPosition = -1;

	constructor(id, kind, xPosition, yPosition) {
		this.id = id;
		this.kind = kind;
		this.xPosition = xPosition;
		this.yPosition = yPosition;
	}

	print() {
		const cell = document.getElementById(newCellId(this.xPosition, this.yPosition));
		if (!cell) {
			throw new Error(`${newCellId(this.xPosition, this.yPosition)} was not found`);
		}
		cell.innerHTML = `<div id="${this.id}" data-kind="${this.kind}" onclick="handleClick('${this.id}')">${this.kind}</div>`;
	}

	static isCorn(kind) {
		return ["corn:0", "corn:1", "corn:2"].includes(kind);
	}
}
