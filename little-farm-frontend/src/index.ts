import Corn from "./entities/Corn";
import Dirt from "./entities/Dirt";
import Grass from "./entities/Grass";
import {
	ActiveElement,
	advanceOneSecond,
	createFarm,
	Farm,
	harvestCorn,
	harvestGrass,
	plantCorn,
	plantWheat,
} from "./services/farm";
import { newCellId } from "./utils";

const MAP_WIDTH = 8;
const MAP_HEIGHT = 8;
const ONE_SECOND = 1000;

let currentFarm: Farm;
let selectedDirtId: string | null = null;
const cells: Map<string, string> = new Map();

async function main() {
	currentFarm = await createFarm(MAP_WIDTH, MAP_HEIGHT);
	createCells(currentFarm);
	populateCells(currentFarm);
	printResources(currentFarm);

	setInterval(async () => {
		currentFarm = await advanceOneSecond(currentFarm.farmId);
		populateCells(currentFarm);
		printResources(currentFarm);
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
			currentFarm = await harvestCorn(currentFarm.farmId, e.id);
			populateCells(currentFarm);
			printResources(currentFarm);
			clearPlantingMenu();
		}).print();
	}
	if (Dirt.isDirt(e.kind)) {
		new Dirt(e.id, e.kind, e.xPosition, e.yPosition, async () => {
			selectedDirtId = e.id;
			printPlantingMenu();
		}).print();
	}
	if (Grass.isGrass(e.kind)) {
		new Grass(e.id, e.kind, e.xPosition, e.yPosition, async () => {
			currentFarm = await harvestGrass(currentFarm.farmId, e.id);
			populateCells(currentFarm);
			printResources(currentFarm);
			clearPlantingMenu();
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

function printPlantingMenu() {
	const menu = document.getElementById("planting-menu");
	if (!menu) {
		throw new Error("planting-menu was not found");
	}
	if (!selectedDirtId) {
		menu.innerHTML = "";
		return;
	}

	const options: string[] = [];
	if (currentFarm.resources.corn >= 1) {
		options.push('<button id="planting-menu__corn">Plant Corn</button>');
	}
	if (currentFarm.resources.seeds >= 1) {
		options.push('<button id="planting-menu__wheat">Plant Wheat</button>');
	}
	menu.innerHTML = options.length > 0 ? options.join("") : "";

	if (currentFarm.resources.corn >= 1) {
		const cornButton = document.getElementById("planting-menu__corn");
		if (cornButton) {
			cornButton.addEventListener("click", handlePlantCorn);
		}
	}
	if (currentFarm.resources.seeds >= 1) {
		const wheatButton = document.getElementById("planting-menu__wheat");
		if (wheatButton) {
			wheatButton.addEventListener("click", handlePlantWheat);
		}
	}
}

async function handlePlantCorn() {
	if (!selectedDirtId) {
		return;
	}
	currentFarm = await plantCorn(currentFarm.farmId, selectedDirtId);
	populateCells(currentFarm);
	printResources(currentFarm);
	clearPlantingMenu();
}

async function handlePlantWheat() {
	if (!selectedDirtId) {
		return;
	}
	currentFarm = await plantWheat(currentFarm.farmId, selectedDirtId);
	populateCells(currentFarm);
	printResources(currentFarm);
	clearPlantingMenu();
}

function clearPlantingMenu() {
	const menu = document.getElementById("planting-menu");
	if (!menu) {
		throw new Error("planting-menu was not found");
	}
	menu.innerHTML = "";
	selectedDirtId = null;
}
