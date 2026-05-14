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
		const cell = document.getElementById(`cell__${this.xPosition}_${this.yPosition}`);
		if (!cell) {
			throw new Error(`cell__${this.xPosition}_${this.yPosition} was not found`);
		}
		cell.innerHTML = `
			<div id="${this.id}" data-kind="grass" onclick="handleClick('${this.id}')">${this.kind}</div>
		`;
	}

	static isGrass(kind) {
		return ["grass:0", "grass:1", "grass:2", "grass:3", "grass:4"].includes(kind);
	}
}
