export default class Dirt {
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
			<div id="${this.id}" data-kind="dirt" onclick="handleClick('${this.id}')">${this.kind}</div>
		`;
	}

	static isDirt(kind) {
		return ["dirt:0"].includes(kind);
	}
}
