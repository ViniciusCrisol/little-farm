import { newCellId } from "../utils";

export default class Dirt {
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

	static isDirt(kind: string) {
		return ["dirt:0"].includes(kind);
	}
}
