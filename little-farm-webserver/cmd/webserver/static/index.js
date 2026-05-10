class Corn {
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
			<div id="${this.id}" onclick="handleClick('${this.id}')">${this.kind}</div>
		`;
	}

	static isCorn(kind) {
		return kind === "corn:0" || kind === "corn:1" || kind === "corn:2";
	}
}

class Dirt {
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
			<div id="${this.id}" onclick="handleClick('${this.id}')">${this.kind}</div>
		`;
	}

	static isDirt(kind) {
		return kind === "grass:0";
	}
}

async function main() {
	const mapWidth = 10;
	const mapHeight = 10;

	const response = await fetch("/api/farms", {
		body: JSON.stringify({
			map_width: mapWidth,
			map_height: mapHeight,
		}),
		method: "POST",
	});
	const jsonResponse = await response.json();

	const farm = {
		farmID: jsonResponse.farm_id,
		mapWidth: jsonResponse.map_width,
		mapHeight: jsonResponse.map_height,
		activeFor: jsonResponse.active_for,
		resources: {
			corn: jsonResponse.resources.corn,
			seeds: jsonResponse.resources.seeds,
			moneyInCents: jsonResponse.resources.money_in_cents,
		},
		activeElements: jsonResponse.active_elements.map((e) => ({
			id: e.id,
			kind: e.kind,
			xPosition: e.x_position,
			yPosition: e.y_position,
		})),
		passiveElements: jsonResponse.passive_elements.map((e) => ({
			id: e.id,
			kind: e.kind,
		})),
		createdAt: jsonResponse.created_at,
		updatedAt: jsonResponse.updated_at,
	};

	create(mapWidth, mapHeight);
	populate(farm);
}
main();

function create(mapWidth, mapHeight) {
	const cells = document.getElementById("cells");
	for (let x = 0; x < mapWidth; x++) {
		for (let y = 0; y < mapHeight; y++) {
			cells.innerHTML += `<div id="cell__${x}_${y}"></div>`;
		}
	}
}

function populate(farm) {
	farm.activeElements.forEach((e) => {
		if (Corn.isCorn(e.kind)) {
			new Corn(e.id, e.kind, e.xPosition, e.yPosition).print();
		}
		if (Dirt.isDirt(e.kind)) {
			new Dirt(e.id, e.kind, e.xPosition, e.yPosition).print();
		}
	});
}
