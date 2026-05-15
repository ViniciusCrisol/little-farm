import { newCellId } from "../utils";

export default class Corn {
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

	static isCorn(kind: string) {
		return ["corn:0", "corn:1", "corn:2"].includes(kind);
	}
}
