import Corn from "./entities/Corn";
import Dirt from "./entities/Dirt";
import Grass from "./entities/Grass";
import { ActiveElement, advanceOneSecond, createFarm, Farm, harvestCorn, harvestGrass } from "./services/farm";
import { newCellId } from "./utils";

const MAP_WIDTH = 8;
const MAP_HEIGHT = 8;
const ONE_SECOND = 1000;

let farmId: string;
const cells: Map<string, string> = new Map();

async function main() {
	const farm = await createFarm(MAP_WIDTH, MAP_HEIGHT);
	createCells(farm);
	populateCells(farm);
	printResources(farm);
	farmId = farm.farmId;

	setInterval(async () => {
		const farm = await advanceOneSecond(farmId);
		populateCells(farm);
		printResources(farm);
	}, ONE_SECOND);
}
main();

function createCells(farm: Farm) {
	const cells = document.getElementById("cells");
	if (!cells) {
		throw new Error("cells was not found");
	}
	for (let x = 0; x < farm.mapWidth; x++) {
		for (let y = 0; y < farm.mapHeight; y++) {
			cells.innerHTML += `<div id="${newCellId(x, y)}"></div>`;
		}
	}
}

function populateCells(farm: Farm) {
	const elementIDs = new Set(farm.activeElements.map((e) => e.id));
	cells.forEach((_, id) => {
		if (!elementIDs.has(id)) {
			removeElement(id);
			cells.delete(id);
		}
	});
	farm.activeElements.forEach((e) => {
		if (cells.has(e.id)) {
			if (cells.get(e.id) !== e.kind) {
				removeElement(e.id);
				cells.set(e.id, e.kind);
				createElement(e);
			}
			return;
		}
		cells.set(e.id, e.kind);
		createElement(e);
	});
}

function createElement(e: ActiveElement) {
	if (Corn.isCorn(e.kind)) {
		new Corn(e.id, e.kind, e.xPosition, e.yPosition, async () => {
			const farm = await harvestCorn(farmId, e.id);
			populateCells(farm);
			printResources(farm);
		}).print();
	}
	if (Dirt.isDirt(e.kind)) {
		new Dirt(e.id, e.kind, e.xPosition, e.yPosition).print();
	}
	if (Grass.isGrass(e.kind)) {
		new Grass(e.id, e.kind, e.xPosition, e.yPosition, async () => {
			const farm = await harvestGrass(farmId, e.id);
			populateCells(farm);
			printResources(farm);
		}).print();
	}
}

function removeElement(elementId: string) {
	document.getElementById(elementId)?.remove();
}

function printResources(farm: Farm) {
	const corn = document.getElementById("resource__corn");
	if (corn) {
		corn.textContent = farm.resources.corn.toString();
	}

	const seeds = document.getElementById("resource__seeds");
	if (seeds) {
		seeds.textContent = farm.resources.seeds.toString();
	}

	const money = document.getElementById("resource__money");
	if (money) {
		money.textContent = (farm.resources.moneyInCents / 100).toFixed(2);
	}
}
