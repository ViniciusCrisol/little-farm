import { createFarm } from "./api/index.js";
import Corn from "./entities/Corn.js";
import Dirt from "./entities/Dirt.js";
import Grass from "./entities/Grass.js";

const state = {};

async function main() {
	const mapWidth = 10;
	const mapHeight = 10;
	const farm = await createFarm(mapWidth, mapHeight);
	createCells(mapWidth, mapHeight);
	populate(farm);
}
main();

function createCells(mapWidth, mapHeight) {
	const cells = document.getElementById("cells");
	for (let x = 0; x < mapWidth; x++) {
		for (let y = 0; y < mapHeight; y++) {
			cells.innerHTML += `<div id="cell__${x}_${y}"></div>`;
		}
	}
}

function populate(farm) {
	const stateElementIDs = Object.keys(state);

	stateElementIDs.forEach((id) => {
		if (!farm.activeElements.find((e) => e.id === id)) {
			document.getElementById(id).remove();
			delete state[id];
		}
	});
	farm.activeElements.forEach((e) => {
		if (stateElementIDs.includes(e.id)) {
			if (state[e.id] !== e.kind) {
				document.getElementById(id).remove();
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
