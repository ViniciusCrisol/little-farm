import { newCellId } from "../utils";

export default class Grass {
	constructor(
		private readonly id: string,
		private readonly kind: string,
		private readonly xPosition: number,
		private readonly yPosition: number,
	) {}

	print() {
		const cell = document.getElementById(newCellId(this.xPosition, this.yPosition));
		if (!cell) {
			throw new Error(`${newCellId(this.xPosition, this.yPosition)} was not found`);
		}
		cell.innerHTML = `<div id="${this.id}" data-kind="${this.kind}" onclick="handleClick('${this.id}')">${this.kind}</div>`;
	}

	static isGrass(kind: string) {
		return ["grass:0", "grass:1", "grass:2", "grass:3", "grass:4"].includes(kind);
	}
}
