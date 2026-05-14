import Corn from "./entities/Corn.js";
import Dirt from "./entities/Dirt.js";
import Grass from "./entities/Grass.js";
import { createFarm } from "./services/farm.js";
import { newCellId } from "./utils.js";

const MAP_WIDTH = 10;
const MAP_HEIGHT = 10;

const state = {};

async function main() {
	const farm = await createFarm(MAP_WIDTH, MAP_HEIGHT);
	createCells(farm);
	populateCells(farm);
}
main();

function createCells(farm) {
	const cells = document.getElementById("cells");
	for (let x = 0; x < farm.mapWidth; x++) {
		for (let y = 0; y < farm.mapHeight; y++) {
			cells.innerHTML += `<div id="${newCellId(x, y)}"></div>`;
		}
	}
}

function populateCells(farm) {
	const printedIDs = new Set(Object.keys(state));
	const elementIDs = new Set(farm.activeElements.map((e) => e.id));
	printedIDs.forEach((id) => {
		if (!elementIDs.has(id)) {
			removeElement(id);
			delete state[id];
		}
	});
	farm.activeElements.forEach((e) => {
		if (printedIDs.has(e.id)) {
			if (state[e.id] !== e.kind) {
				removeElement(e.id);
				state[e.id] = e.kind;
				createElement(e);
			}
			return;
		}
		state[e.id] = e.kind;
		createElement(e);
	});
}

function createElement(e) {
	if (Corn.isCorn(e.kind)) {
		new Corn(e.id, e.kind, e.xPosition, e.yPosition).print();
	}
	if (Dirt.isDirt(e.kind)) {
		new Dirt(e.id, e.kind, e.xPosition, e.yPosition).print();
	}
	if (Grass.isGrass(e.kind)) {
		new Grass(e.id, e.kind, e.xPosition, e.yPosition).print();
	}
}

function removeElement(elementId) {
	document.getElementById(elementId)?.remove();
}
