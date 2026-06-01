type FarmResponse = {
	farm_id: string;
	map_width: number;
	map_height: number;
	active_for: number;
	resources: {
		corn: number;
		seeds: number;
		money_in_cents: number;
	};
	active_elements: Array<{
		id: string;
		kind: string;
		x_position: number;
		y_position: number;
	}>;
	passive_elements: Array<{
		id: string;
		kind: string;
	}>;
	created_at: string;
	updated_at: string;
};

export type Farm = {
	farmId: string;
	mapWidth: number;
	mapHeight: number;
	activeFor: number;
	resources: Resources;
	activeElements: Array<ActiveElement>;
	passiveElements: Array<PassiveElement>;
	createdAt: Date;
	updatedAt: Date;
};

export type Resources = {
	corn: number;
	seeds: number;
	moneyInCents: number;
};

export type ActiveElement = {
	id: string;
	kind: string;
	xPosition: number;
	yPosition: number;
};

export type PassiveElement = {
	id: string;
	kind: string;
};

function farmResponseToFarm(response: FarmResponse): Farm {
	return {
		farmId: response.farm_id,
		mapWidth: response.map_width,
		mapHeight: response.map_height,
		activeFor: response.active_for,
		resources: {
			corn: response.resources.corn,
			seeds: response.resources.seeds,
			moneyInCents: response.resources.money_in_cents,
		},
		activeElements: response.active_elements.map((e) => ({
			id: e.id,
			kind: e.kind,
			xPosition: e.x_position,
			yPosition: e.y_position,
		})),
		passiveElements: response.passive_elements.map((e) => ({
			id: e.id,
			kind: e.kind,
		})),
		createdAt: new Date(response.created_at),
		updatedAt: new Date(response.updated_at),
	};
}

export async function createFarm(mapWidth: number, mapHeight: number): Promise<Farm> {
	const response = await fetch("/api/farms", {
		body: JSON.stringify({
			map_width: mapWidth,
			map_height: mapHeight,
		}),
		method: "POST",
	});
	return farmResponseToFarm(await response.json());
}

export async function advanceOneSecond(farmId: string): Promise<Farm> {
	const response = await fetch(`/api/farms/${farmId}/advance-one-second`, {
		method: "POST",
	});
	return farmResponseToFarm(await response.json());
}

export async function plantCorn(farmId: string, dirtId: string): Promise<Farm> {
	const response = await fetch(`/api/farms/${farmId}/dirt/${dirtId}/plant-corn`, {
		method: "POST",
	});
	return farmResponseToFarm(await response.json());
}

export async function plantWheat(farmId: string, dirtId: string): Promise<Farm> {
	const response = await fetch(`/api/farms/${farmId}/dirt/${dirtId}/plant-wheat`, {
		method: "POST",
	});
	return farmResponseToFarm(await response.json());
}

export async function harvestCorn(farmId: string, cornId: string): Promise<Farm> {
	const response = await fetch(`/api/farms/${farmId}/corn/${cornId}/harvest`, {
		method: "POST",
	});
	return farmResponseToFarm(await response.json());
}

export async function harvestGrass(farmId: string, grassId: string): Promise<Farm> {
	const response = await fetch(`/api/farms/${farmId}/grass/${grassId}/harvest`, {
		method: "POST",
	});
	return farmResponseToFarm(await response.json());
}
