import { newCellId } from "../utils.js";

export default class Grass {
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

	static isGrass(kind) {
		return ["grass:0", "grass:1", "grass:2", "grass:3", "grass:4"].includes(kind);
	}
}
