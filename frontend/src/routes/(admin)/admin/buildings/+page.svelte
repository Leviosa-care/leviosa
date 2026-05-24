<script lang="ts">
	import type { PageData } from "./$types";
	import type { Building, Room, Allocation, Partner } from "./+page.server";
	import Modal from "$lib/ui/Modal.svelte";
	import Drawer from "$lib/ui/Drawer.svelte";
	import {
		Building2,
		ChevronDown,
		ChevronRight,
		Plus,
		Pencil,
		Users,
		DoorOpen,
		CalendarRange,
		Link2,
		PowerOff,
		AlertCircle,
		CheckCircle,
		X,
	} from "@lucide/svelte";

	let { data }: { data: PageData } = $props();

	// --- Types ---
	interface RoomWithAllocations extends Room {
		allocations: Allocation[];
		allocationsLoaded: boolean;
		allocationsExpanded: boolean;
	}

	interface BuildingWithRooms extends Building {
		rooms: RoomWithAllocations[];
		roomsLoaded: boolean;
		roomsExpanded: boolean;
	}

	// --- State ---
	let buildings = $state<BuildingWithRooms[]>(
		data.buildings.map((b) => ({
			...b,
			rooms: [],
			roomsLoaded: false,
			roomsExpanded: false,
		}))
	);
	let partners = $state<Partner[]>(data.partners);
	let isMobile = $state(false);
	let loading = $state<Record<string, boolean>>({});
	let feedback = $state<{ type: "success" | "error"; message: string } | null>(null);

	// Modal states
	let createBuildingOpen = $state(false);
	let editBuildingOpen = $state(false);
	let createRoomOpen = $state(false);
	let editRoomOpen = $state(false);
	let createAllocationOpen = $state(false);
	let editPeriodOpen = $state(false);
	let deactivateConfirmOpen = $state(false);

	let selectedBuilding = $state<BuildingWithRooms | null>(null);
	let selectedRoom = $state<RoomWithAllocations | null>(null);
	let selectedAllocation = $state<Allocation | null>(null);

	// Form state
	let newBuilding = $state({ name: "", address: "", city: "", postal_code: "", country: "France", phone: "", email: "" });
	let editBuildingForm = $state({ name: "", address: "", city: "", postal_code: "", country: "", phone: "", email: "" });
	let newRoom = $state({ name: "", capacity: 1 });
	let editRoomForm = $state({ name: "", capacity: 1 });
	let newAllocationType = $state<"shared" | "dedicated">("shared");
	let newAllocationPartnerId = $state("");
	let newAllocationStartDate = $state("");
	let newAllocationEndDate = $state("");
	let editPeriodStartDate = $state("");
	let editPeriodEndDate = $state("");

	// --- Mobile detection ---
	$effect(() => {
		isMobile = window.innerWidth < 768;
		const onResize = () => { isMobile = window.innerWidth < 768; };
		window.addEventListener("resize", onResize);
		return () => window.removeEventListener("resize", onResize);
	});

	// --- Feedback ---
	function showFeedback(type: "success" | "error", message: string) {
		feedback = { type, message };
		setTimeout(() => { feedback = null; }, 4000);
	}

	// --- Helpers ---
	function partnerName(userId: string): string {
		const p = partners.find((p) => p.user_id === userId);
		return (p?.user_id ?? userId).slice(0, 8) + "…";
	}

	function formatDate(d: string | null | undefined): string {
		if (!d) return "—";
		return new Date(d).toLocaleDateString("fr-FR", { day: "numeric", month: "short", year: "numeric" });
	}

	function allocationStatusLabel(a: Allocation): string {
		if (!a.is_active) return "Inactive";
		if (a.allocation_type === "shared") return "Active";
		if (a.end_date && new Date(a.end_date) < new Date()) return "Expirée";
		if (a.start_date && new Date(a.start_date) > new Date()) return "Planifiée";
		return "Active";
	}

	function allocationStatusColor(a: Allocation): string {
		const label = allocationStatusLabel(a);
		if (label === "Active") return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300";
		if (label === "Planifiée") return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300";
		if (label === "Expirée") return "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300";
		return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400";
	}

	function resetNewBuilding() {
		newBuilding = { name: "", address: "", city: "", postal_code: "", country: "France", phone: "", email: "" };
	}

	function resetNewRoom() {
		newRoom = { name: "", capacity: 1 };
	}

	function resetNewAllocation() {
		newAllocationType = "shared";
		newAllocationPartnerId = "";
		newAllocationStartDate = "";
		newAllocationEndDate = "";
	}

	// --- API calls ---
	async function apiFetch(path: string, method: string = "GET", body?: unknown): Promise<Response> {
		const init: RequestInit = {
			method,
			headers: { "Content-Type": "application/json" },
		};
		if (body !== undefined) init.body = JSON.stringify(body);
		return fetch(path, init);
	}

	// --- Building actions ---
	async function submitCreateBuilding() {
		if (!newBuilding.name || !newBuilding.address || !newBuilding.city || !newBuilding.postal_code) {
			showFeedback("error", "Veuillez remplir tous les champs obligatoires");
			return;
		}
		loading["create-building"] = true;
		try {
			const res = await apiFetch("/buildings", "POST", { ...newBuilding, is_active: true });
			if (!res.ok) throw new Error(await res.text());
			const created: Building = await res.json();
			buildings.unshift({ ...created, rooms: [], roomsLoaded: false, roomsExpanded: false });
			createBuildingOpen = false;
			resetNewBuilding();
			showFeedback("success", "Bâtiment créé avec succès");
		} catch (e) {
			showFeedback("error", "Erreur lors de la création du bâtiment");
		} finally {
			loading["create-building"] = false;
		}
	}

	async function submitEditBuilding() {
		if (!selectedBuilding) return;
		loading["edit-building"] = true;
		try {
			const res = await apiFetch(`/buildings/${selectedBuilding.id}`, "PUT", {
				id: selectedBuilding.id,
				name: editBuildingForm.name,
				address: editBuildingForm.address,
				city: editBuildingForm.city,
				postal_code: editBuildingForm.postal_code,
				country: editBuildingForm.country,
				phone: editBuildingForm.phone || undefined,
				email: editBuildingForm.email || undefined,
			});
			if (!res.ok) throw new Error(await res.text());
			const updated: Building = await res.json();
			const idx = buildings.findIndex((b) => b.id === updated.id);
			if (idx >= 0) {
				buildings[idx] = { ...buildings[idx], ...updated };
			}
			editBuildingOpen = false;
			showFeedback("success", "Bâtiment modifié avec succès");
		} catch {
			showFeedback("error", "Erreur lors de la modification du bâtiment");
		} finally {
			loading["edit-building"] = false;
		}
	}

	// --- Room actions ---
	async function toggleBuildingRooms(b: BuildingWithRooms) {
		b.roomsExpanded = !b.roomsExpanded;
		if (!b.roomsLoaded) {
			await loadRooms(b);
		}
	}

	async function loadRooms(b: BuildingWithRooms) {
		loading[`rooms-${b.id}`] = true;
		try {
			const res = await apiFetch(`/buildings/${b.id}/rooms`, "GET");
			if (!res.ok) throw new Error("fetch failed");
			const rooms: Room[] = await res.json();
			b.rooms = rooms.map((r) => ({ ...r, allocations: [], allocationsLoaded: false, allocationsExpanded: false }));
			b.roomsLoaded = true;
		} catch {
			showFeedback("error", "Erreur lors du chargement des cabinets");
		} finally {
			loading[`rooms-${b.id}`] = false;
		}
	}

	async function submitCreateRoom() {
		if (!selectedBuilding || !newRoom.name) {
			showFeedback("error", "Veuillez remplir le nom du cabinet");
			return;
		}
		loading["create-room"] = true;
		try {
			const res = await apiFetch("/rooms", "POST", {
				building_id: selectedBuilding.id,
				name: newRoom.name,
				capacity: newRoom.capacity,
				is_active: true,
			});
			if (!res.ok) throw new Error(await res.text());
			const created: Room = await res.json();
			if (selectedBuilding.roomsLoaded) {
				selectedBuilding.rooms.push({ ...created, allocations: [], allocationsLoaded: false, allocationsExpanded: false });
			}
			createRoomOpen = false;
			resetNewRoom();
			showFeedback("success", "Cabinet créé avec succès");
		} catch {
			showFeedback("error", "Erreur lors de la création du cabinet");
		} finally {
			loading["create-room"] = false;
		}
	}

	async function submitEditRoom() {
		if (!selectedRoom) return;
		loading["edit-room"] = true;
		try {
			const res = await apiFetch(`/rooms/${selectedRoom.id}`, "PUT", {
				id: selectedRoom.id,
				name: editRoomForm.name,
				capacity: editRoomForm.capacity,
			});
			if (!res.ok) throw new Error(await res.text());
			const updated: Room = await res.json();
			// Update in parent building
			for (const b of buildings) {
				const idx = b.rooms.findIndex((r) => r.id === updated.id);
				if (idx >= 0) {
					b.rooms[idx] = { ...b.rooms[idx], ...updated };
					break;
				}
			}
			editRoomOpen = false;
			showFeedback("success", "Cabinet modifié avec succès");
		} catch {
			showFeedback("error", "Erreur lors de la modification du cabinet");
		} finally {
			loading["edit-room"] = false;
		}
	}

	// --- Allocation actions ---
	async function toggleRoomAllocations(b: BuildingWithRooms, r: RoomWithAllocations) {
		r.allocationsExpanded = !r.allocationsExpanded;
		if (!r.allocationsLoaded) {
			await loadAllocations(b, r);
		}
	}

	async function loadAllocations(b: BuildingWithRooms, r: RoomWithAllocations) {
		loading[`alloc-${r.id}`] = true;
		try {
			const res = await apiFetch(`/rooms/${r.id}/allocations?active_only=false`, "GET");
			if (!res.ok) throw new Error("fetch failed");
			const allocations: Allocation[] = await res.json();
			r.allocations = allocations;
			r.allocationsLoaded = true;
		} catch {
			showFeedback("error", "Erreur lors du chargement des allocations");
		} finally {
			loading[`alloc-${r.id}`] = false;
		}
	}

	async function submitCreateAllocation() {
		if (!selectedRoom || !newAllocationPartnerId) {
			showFeedback("error", "Veuillez sélectionner un partenaire");
			return;
		}
		loading["create-allocation"] = true;
		try {
			let res: Response;
			if (newAllocationType === "shared") {
				res = await apiFetch("/allocations/shared", "POST", {
					room_id: selectedRoom.id,
					user_id: newAllocationPartnerId,
				});
			} else {
				if (!newAllocationStartDate) {
					showFeedback("error", "Veuillez renseigner la date de début");
					loading["create-allocation"] = false;
					return;
				}
				const body: Record<string, string> = {
					room_id: selectedRoom.id,
					user_id: newAllocationPartnerId,
					start_date: newAllocationStartDate + "T00:00:00Z",
				};
				if (newAllocationEndDate) body.end_date = newAllocationEndDate + "T00:00:00Z";
				res = await apiFetch("/allocations/dedicated", "POST", body);
			}
			if (!res.ok) throw new Error(await res.text());
			// Reload allocations for this room
			const parentBuilding = buildings.find((b) => b.rooms.some((r) => r.id === selectedRoom!.id));
			if (parentBuilding) {
				const room = parentBuilding.rooms.find((r) => r.id === selectedRoom!.id);
				if (room) {
					room.allocationsLoaded = false;
					await loadAllocations(parentBuilding, room);
				}
			}
			createAllocationOpen = false;
			resetNewAllocation();
			showFeedback("success", "Allocation créée avec succès");
		} catch {
			showFeedback("error", "Erreur lors de la création de l'allocation");
		} finally {
			loading["create-allocation"] = false;
		}
	}

	async function submitDeactivate() {
		if (!selectedAllocation) return;
		loading["deactivate"] = true;
		try {
			const res = await apiFetch(`/allocations/${selectedAllocation.id}/deactivate`, "POST");
			if (!res.ok) throw new Error("fetch failed");
			// Update locally
			outer: for (const b of buildings) {
				for (const r of b.rooms) {
					const a = r.allocations.find((a) => a.id === selectedAllocation!.id);
					if (a) { a.is_active = false; break outer; }
				}
			}
			deactivateConfirmOpen = false;
			showFeedback("success", "Allocation désactivée");
		} catch {
			showFeedback("error", "Erreur lors de la désactivation");
		} finally {
			loading["deactivate"] = false;
		}
	}

	async function submitEditPeriod() {
		if (!selectedAllocation) return;
		loading["edit-period"] = true;
		try {
			const body: Record<string, string> = {
				id: selectedAllocation.id,
				start_date: editPeriodStartDate + "T00:00:00Z",
			};
			if (editPeriodEndDate) body.end_date = editPeriodEndDate + "T00:00:00Z";
			const res = await apiFetch(`/allocations/${selectedAllocation.id}/period`, "PUT", body);
			if (!res.ok) throw new Error(await res.text());
			// Reload allocations
			for (const b of buildings) {
				for (const r of b.rooms) {
					if (r.allocations.some((a) => a.id === selectedAllocation!.id)) {
						r.allocationsLoaded = false;
						await loadAllocations(b, r);
						break;
					}
				}
			}
			editPeriodOpen = false;
			showFeedback("success", "Période modifiée avec succès");
		} catch {
			showFeedback("error", "Erreur lors de la modification de la période");
		} finally {
			loading["edit-period"] = false;
		}
	}

	// --- Open modal helpers ---
	function openCreateBuilding() {
		resetNewBuilding();
		createBuildingOpen = true;
	}

	function openEditBuilding(b: BuildingWithRooms) {
		selectedBuilding = b;
		editBuildingForm = {
			name: b.name,
			address: b.address,
			city: b.city,
			postal_code: b.postal_code,
			country: b.country,
			phone: b.phone,
			email: b.email,
		};
		editBuildingOpen = true;
	}

	function openCreateRoom(b: BuildingWithRooms) {
		selectedBuilding = b;
		resetNewRoom();
		createRoomOpen = true;
	}

	function openEditRoom(b: BuildingWithRooms, r: RoomWithAllocations) {
		selectedBuilding = b;
		selectedRoom = r;
		editRoomForm = { name: r.name, capacity: r.capacity };
		editRoomOpen = true;
	}

	function openCreateAllocation(b: BuildingWithRooms, r: RoomWithAllocations) {
		selectedBuilding = b;
		selectedRoom = r;
		resetNewAllocation();
		createAllocationOpen = true;
	}

	function openDeactivate(a: Allocation) {
		selectedAllocation = a;
		deactivateConfirmOpen = true;
	}

	function openEditPeriod(a: Allocation) {
		selectedAllocation = a;
		editPeriodStartDate = a.start_date ? a.start_date.split("T")[0] : "";
		editPeriodEndDate = a.end_date ? a.end_date.split("T")[0] : "";
		editPeriodOpen = true;
	}

	// Filter verified partners for allocation
	let verifiedPartners = $derived(
		partners.filter((p) => p.stripe_onboarding_complete && p.stripe_account_status === "active")
	);
</script>

<svelte:head>
	<title>Bâtiments &amp; Cabinets | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="text-3xl lg:text-4xl font-bold mb-2 text-foreground">Bâtiments &amp; Cabinets</h1>
			<p class="text-muted-foreground">Gérer les bâtiments, cabinets et allocations</p>
		</div>
		<button
			onclick={openCreateBuilding}
			class="inline-flex items-center gap-2 px-4 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm"
		>
			<Plus size={16} />
			Nouveau bâtiment
		</button>
	</div>

	{#if feedback}
		<div class="mb-6 flex items-center gap-3 px-4 py-3 rounded-lg border {feedback.type === 'success'
			? 'bg-green-50 border-green-200 text-green-800'
			: 'bg-red-50 border-red-200 text-red-800'}">
			{#if feedback.type === 'success'}
				<CheckCircle size={20} />
			{:else}
				<AlertCircle size={20} />
			{/if}
			<p class="text-sm font-medium">{feedback.message}</p>
		</div>
	{/if}

	{#if buildings.length === 0}
		<div class="text-center py-16">
			<Building2 size={48} class="mx-auto text-muted-foreground mb-4" />
			<p class="text-muted-foreground text-lg">Aucun bâtiment enregistré</p>
			<p class="text-muted-foreground text-sm mt-1">Créez votre premier bâtiment pour commencer</p>
		</div>
	{:else}
		<div class="space-y-4">
			{#each buildings as building (building.id)}
				<div class="bg-background border border-border-card rounded-lg overflow-hidden">
					<!-- Building header -->
					<div
						role="button"
						tabindex="0"
						onclick={() => toggleBuildingRooms(building)}
						onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); toggleBuildingRooms(building); } }}
						class="w-full flex items-center justify-between px-6 py-4 hover:bg-muted/30 transition-colors text-left cursor-pointer"
					>
						<div class="flex items-center gap-3">
							{#if building.roomsExpanded}
								<ChevronDown size={20} class="text-muted-foreground" />
							{:else}
								<ChevronRight size={20} class="text-muted-foreground" />
							{/if}
							<div class="w-10 h-10 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center">
								<Building2 size={20} class="text-blue-600 dark:text-blue-400" />
							</div>
							<div>
								<h3 class="text-base font-semibold text-foreground">{building.name}</h3>
								<p class="text-sm text-muted-foreground">
									{building.city}{building.address ? ` — ${building.address}` : ''}
									{#if building.roomsLoaded}
										· {building.rooms.length} cabinet{building.rooms.length !== 1 ? 's' : ''}
									{/if}
								</p>
							</div>
						</div>
						<div class="flex items-center gap-2" onclick={(e) => e.stopPropagation()}>
							<button
								onclick={() => openEditBuilding(building)}
								class="p-2 rounded-md text-foreground-alt hover:text-foreground hover:bg-muted transition-colors"
								title="Modifier"
							>
								<Pencil size={16} />
							</button>
							<button
								onclick={() => openCreateRoom(building)}
								class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-foreground-alt hover:text-foreground bg-muted hover:bg-muted/80 rounded-md transition-colors"
								title="Ajouter un cabinet"
							>
								<Plus size={14} />
								Cabinet
							</button>
						</div>
					</div>

					<!-- Rooms section -->
					{#if building.roomsExpanded}
						<div class="border-t border-border-card">
							{#if loading[`rooms-${building.id}`]}
								<div class="px-6 py-8 text-center text-muted-foreground text-sm">Chargement des cabinets…</div>
							{:else if building.rooms.length === 0}
								<div class="px-6 py-8 text-center text-muted-foreground text-sm">Aucun cabinet dans ce bâtiment</div>
							{:else}
								<div class="divide-y divide-border-card">
									{#each building.rooms as room (room.id)}
										<div>
											<!-- Room header -->
											<div
												onclick={() => toggleRoomAllocations(building, room)}
												role="button"
												tabindex="0"
												class="w-full flex items-center justify-between px-6 py-3 pl-12 hover:bg-muted/30 transition-colors text-left"
												onkeydown={(e: KeyboardEvent) => { if (e.key === "Enter" || e.key === " ") { e.preventDefault(); toggleRoomAllocations(building, room); } }}
											>
												<div class="flex items-center gap-3">
													{#if room.allocationsExpanded}
														<ChevronDown size={18} class="text-muted-foreground" />
													{:else}
														<ChevronRight size={18} class="text-muted-foreground" />
													{/if}
													<div class="w-8 h-8 rounded-md bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
														<DoorOpen size={16} class="text-green-600 dark:text-green-400" />
													</div>
													<div>
														<span class="text-sm font-medium text-foreground">{room.name}</span>
														<span class="ml-2 text-xs text-muted-foreground">Capacité: {room.capacity}</span>
														{#if room.room_number}
															<span class="ml-2 text-xs text-muted-foreground">N°{room.room_number}</span>
														{/if}
													</div>
												</div>
												<div class="flex items-center gap-2" onclick={(e) => e.stopPropagation()}>
													<button
														onclick={() => openEditRoom(building, room)}
														class="p-1.5 rounded-md text-foreground-alt hover:text-foreground hover:bg-muted transition-colors"
														title="Modifier"
													>
														<Pencil size={14} />
													</button>
													<button
														onclick={() => openCreateAllocation(building, room)}
														class="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium text-foreground-alt hover:text-foreground bg-muted hover:bg-muted/80 rounded-md transition-colors"
														title="Gérer les allocations"
													>
														<Users size={12} />
														Allocations
													</button>
												</div>
											</div>

											<!-- Allocations section -->
											{#if room.allocationsExpanded}
												<div class="border-t border-border-card bg-muted/10">
													{#if loading[`alloc-${room.id}`]}
														<div class="px-6 py-4 pl-20 text-center text-muted-foreground text-xs">Chargement…</div>
													{:else if room.allocations.length === 0}
														<div class="px-6 py-4 pl-20 text-center text-muted-foreground text-xs">Aucune allocation</div>
													{:else}
														<div class="px-6 py-3 pl-20 space-y-2">
															{#each room.allocations as alloc (alloc.id)}
																<div class="flex items-center justify-between py-2 px-3 bg-background border border-border-card rounded-md">
																	<div class="flex items-center gap-3">
																		{#if alloc.allocation_type === "shared"}
																			<Link2 size={14} class="text-blue-500" />
																			<span class="text-xs font-medium text-foreground">Partagée</span>
																		{:else}
																			<CalendarRange size={14} class="text-purple-500" />
																			<span class="text-xs font-medium text-foreground">Dédiée</span>
																			<span class="text-xs text-muted-foreground">
																				{formatDate(alloc.start_date)}{alloc.end_date ? ` → ${formatDate(alloc.end_date)}` : ''}
																			</span>
																		{/if}
																		<span class="text-xs text-muted-foreground">
																			{partnerName(alloc.user_id)}
																		</span>
																	</div>
																	<div class="flex items-center gap-2">
																		<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium {allocationStatusColor(alloc)}">
																			{allocationStatusLabel(alloc)}
																		</span>
																		{#if alloc.is_active && alloc.allocation_type === "dedicated"}
																			<button
																				onclick={() => openEditPeriod(alloc)}
																				class="p-1 rounded text-foreground-alt hover:text-foreground hover:bg-muted transition-colors"
																				title="Modifier la période"
																			>
																				<Pencil size={12} />
																			</button>
																		{/if}
																		{#if alloc.is_active}
																			<button
																				onclick={() => openDeactivate(alloc)}
																				class="p-1 rounded text-red-500 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
																				title="Désactiver"
																			>
																				<PowerOff size={12} />
																			</button>
																		{/if}
																	</div>
																</div>
															{/each}
														</div>
													{/if}
												</div>
											{/if}
										</div>
									{/each}
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- ===================== MODALS ===================== -->

<!-- Create Building Modal -->
{#if createBuildingOpen}
	{#if isMobile}
		<Drawer bind:isOpen={createBuildingOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">Nouveau bâtiment</h2>
					<button type="button" onclick={() => (createBuildingOpen = false)} class="p-2 hover:bg-muted rounded-md transition-colors">
						<X class="text-foreground size-5" />
					</button>
				</div>
			</div>
			<div class="pt-4 space-y-4">
				<div>
					<label for="cb-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="cb-name" type="text" bind:value={newBuilding.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cb-address" class="block text-sm font-medium text-foreground mb-1">Adresse *</label>
					<input id="cb-address" type="text" bind:value={newBuilding.address} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="cb-city" class="block text-sm font-medium text-foreground mb-1">Ville *</label>
						<input id="cb-city" type="text" bind:value={newBuilding.city} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
					<div>
						<label for="cb-postal" class="block text-sm font-medium text-foreground mb-1">Code postal *</label>
						<input id="cb-postal" type="text" bind:value={newBuilding.postal_code} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
				</div>
				<div>
					<label for="cb-country" class="block text-sm font-medium text-foreground mb-1">Pays</label>
					<input id="cb-country" type="text" bind:value={newBuilding.country} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cb-phone" class="block text-sm font-medium text-foreground mb-1">Téléphone</label>
					<input id="cb-phone" type="text" bind:value={newBuilding.phone} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cb-email" class="block text-sm font-medium text-foreground mb-1">Email</label>
					<input id="cb-email" type="email" bind:value={newBuilding.email} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (createBuildingOpen = false)} class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm">Annuler</button>
					<button onclick={submitCreateBuilding} disabled={loading["create-building"]} class="px-4 py-2 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm disabled:opacity-50">
						{loading["create-building"] ? "Création…" : "Créer"}
					</button>
				</div>
			</div>
		</Drawer>
	{:else}
		<Modal bind:isOpen={createBuildingOpen} title="Nouveau bâtiment">
			<div class="mt-4 space-y-4">
				<div>
					<label for="cb-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="cb-name" type="text" bind:value={newBuilding.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cb-address" class="block text-sm font-medium text-foreground mb-1">Adresse *</label>
					<input id="cb-address" type="text" bind:value={newBuilding.address} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="cb-city" class="block text-sm font-medium text-foreground mb-1">Ville *</label>
						<input id="cb-city" type="text" bind:value={newBuilding.city} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
					<div>
						<label for="cb-postal" class="block text-sm font-medium text-foreground mb-1">Code postal *</label>
						<input id="cb-postal" type="text" bind:value={newBuilding.postal_code} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
				</div>
				<div>
					<label for="cb-country" class="block text-sm font-medium text-foreground mb-1">Pays</label>
					<input id="cb-country" type="text" bind:value={newBuilding.country} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cb-phone" class="block text-sm font-medium text-foreground mb-1">Téléphone</label>
					<input id="cb-phone" type="text" bind:value={newBuilding.phone} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cb-email" class="block text-sm font-medium text-foreground mb-1">Email</label>
					<input id="cb-email" type="email" bind:value={newBuilding.email} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (createBuildingOpen = false)} class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium">Annuler</button>
					<button onclick={submitCreateBuilding} disabled={loading["create-building"]} class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium disabled:opacity-50">
						{loading["create-building"] ? "Création…" : "Créer"}
					</button>
				</div>
			</div>
		</Modal>
	{/if}
{/if}

<!-- Edit Building Modal -->
{#if editBuildingOpen && selectedBuilding}
	{#if isMobile}
		<Drawer bind:isOpen={editBuildingOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">Modifier le bâtiment</h2>
					<button type="button" onclick={() => (editBuildingOpen = false)} class="p-2 hover:bg-muted rounded-md transition-colors">
						<X class="text-foreground size-5" />
					</button>
				</div>
			</div>
			<div class="pt-4 space-y-4">
				<div>
					<label for="eb-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="eb-name" type="text" bind:value={editBuildingForm.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="eb-address" class="block text-sm font-medium text-foreground mb-1">Adresse *</label>
					<input id="eb-address" type="text" bind:value={editBuildingForm.address} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="eb-city" class="block text-sm font-medium text-foreground mb-1">Ville *</label>
						<input id="eb-city" type="text" bind:value={editBuildingForm.city} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
					<div>
						<label for="eb-postal" class="block text-sm font-medium text-foreground mb-1">Code postal *</label>
						<input id="eb-postal" type="text" bind:value={editBuildingForm.postal_code} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
				</div>
				<div>
					<label for="eb-country" class="block text-sm font-medium text-foreground mb-1">Pays</label>
					<input id="eb-country" type="text" bind:value={editBuildingForm.country} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="eb-phone" class="block text-sm font-medium text-foreground mb-1">Téléphone</label>
					<input id="eb-phone" type="text" bind:value={editBuildingForm.phone} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="eb-email" class="block text-sm font-medium text-foreground mb-1">Email</label>
					<input id="eb-email" type="email" bind:value={editBuildingForm.email} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (editBuildingOpen = false)} class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm">Annuler</button>
					<button onclick={submitEditBuilding} disabled={loading["edit-building"]} class="px-4 py-2 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm disabled:opacity-50">
						{loading["edit-building"] ? "Enregistrement…" : "Enregistrer"}
					</button>
				</div>
			</div>
		</Drawer>
	{:else}
		<Modal bind:isOpen={editBuildingOpen} title="Modifier le bâtiment">
			<div class="mt-4 space-y-4">
				<div>
					<label for="eb-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="eb-name" type="text" bind:value={editBuildingForm.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="eb-address" class="block text-sm font-medium text-foreground mb-1">Adresse *</label>
					<input id="eb-address" type="text" bind:value={editBuildingForm.address} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="eb-city" class="block text-sm font-medium text-foreground mb-1">Ville *</label>
						<input id="eb-city" type="text" bind:value={editBuildingForm.city} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
					<div>
						<label for="eb-postal" class="block text-sm font-medium text-foreground mb-1">Code postal *</label>
						<input id="eb-postal" type="text" bind:value={editBuildingForm.postal_code} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
				</div>
				<div>
					<label for="eb-country" class="block text-sm font-medium text-foreground mb-1">Pays</label>
					<input id="eb-country" type="text" bind:value={editBuildingForm.country} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="eb-phone" class="block text-sm font-medium text-foreground mb-1">Téléphone</label>
					<input id="eb-phone" type="text" bind:value={editBuildingForm.phone} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="eb-email" class="block text-sm font-medium text-foreground mb-1">Email</label>
					<input id="eb-email" type="email" bind:value={editBuildingForm.email} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (editBuildingOpen = false)} class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium">Annuler</button>
					<button onclick={submitEditBuilding} disabled={loading["edit-building"]} class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium disabled:opacity-50">
						{loading["edit-building"] ? "Enregistrement…" : "Enregistrer"}
					</button>
				</div>
			</div>
		</Modal>
	{/if}
{/if}

<!-- Create Room Modal -->
{#if createRoomOpen && selectedBuilding}
	{#if isMobile}
		<Drawer bind:isOpen={createRoomOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">Nouveau cabinet</h2>
					<button type="button" onclick={() => (createRoomOpen = false)} class="p-2 hover:bg-muted rounded-md transition-colors">
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">Dans « {selectedBuilding.name} »</p>
			</div>
			<div class="pt-4 space-y-4">
				<div>
					<label for="cr-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="cr-name" type="text" bind:value={newRoom.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cr-cap" class="block text-sm font-medium text-foreground mb-1">Capacité *</label>
					<input id="cr-cap" type="number" min="1" max="50" bind:value={newRoom.capacity} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (createRoomOpen = false)} class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm">Annuler</button>
					<button onclick={submitCreateRoom} disabled={loading["create-room"]} class="px-4 py-2 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm disabled:opacity-50">
						{loading["create-room"] ? "Création…" : "Créer"}
					</button>
				</div>
			</div>
		</Drawer>
	{:else}
		<Modal bind:isOpen={createRoomOpen} title="Nouveau cabinet" description="Dans « {selectedBuilding.name} »">
			<div class="mt-4 space-y-4">
				<div>
					<label for="cr-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="cr-name" type="text" bind:value={newRoom.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="cr-cap" class="block text-sm font-medium text-foreground mb-1">Capacité *</label>
					<input id="cr-cap" type="number" min="1" max="50" bind:value={newRoom.capacity} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (createRoomOpen = false)} class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium">Annuler</button>
					<button onclick={submitCreateRoom} disabled={loading["create-room"]} class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium disabled:opacity-50">
						{loading["create-room"] ? "Création…" : "Créer"}
					</button>
				</div>
			</div>
		</Modal>
	{/if}
{/if}

<!-- Edit Room Modal -->
{#if editRoomOpen && selectedRoom}
	{#if isMobile}
		<Drawer bind:isOpen={editRoomOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">Modifier le cabinet</h2>
					<button type="button" onclick={() => (editRoomOpen = false)} class="p-2 hover:bg-muted rounded-md transition-colors">
						<X class="text-foreground size-5" />
					</button>
				</div>
			</div>
			<div class="pt-4 space-y-4">
				<div>
					<label for="er-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="er-name" type="text" bind:value={editRoomForm.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="er-cap" class="block text-sm font-medium text-foreground mb-1">Capacité *</label>
					<input id="er-cap" type="number" min="1" max="50" bind:value={editRoomForm.capacity} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (editRoomOpen = false)} class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm">Annuler</button>
					<button onclick={submitEditRoom} disabled={loading["edit-room"]} class="px-4 py-2 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm disabled:opacity-50">
						{loading["edit-room"] ? "Enregistrement…" : "Enregistrer"}
					</button>
				</div>
			</div>
		</Drawer>
	{:else}
		<Modal bind:isOpen={editRoomOpen} title="Modifier le cabinet">
			<div class="mt-4 space-y-4">
				<div>
					<label for="er-name" class="block text-sm font-medium text-foreground mb-1">Nom *</label>
					<input id="er-name" type="text" bind:value={editRoomForm.name} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="er-cap" class="block text-sm font-medium text-foreground mb-1">Capacité *</label>
					<input id="er-cap" type="number" min="1" max="50" bind:value={editRoomForm.capacity} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (editRoomOpen = false)} class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium">Annuler</button>
					<button onclick={submitEditRoom} disabled={loading["edit-room"]} class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium disabled:opacity-50">
						{loading["edit-room"] ? "Enregistrement…" : "Enregistrer"}
					</button>
				</div>
			</div>
		</Modal>
	{/if}
{/if}

<!-- Create Allocation Modal -->
{#if createAllocationOpen && selectedRoom}
	{#if isMobile}
		<Drawer bind:isOpen={createAllocationOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">Nouvelle allocation</h2>
					<button type="button" onclick={() => (createAllocationOpen = false)} class="p-2 hover:bg-muted rounded-md transition-colors">
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">Cabinet « {selectedRoom.name} »</p>
			</div>
			<div class="pt-4 space-y-4">
				<!-- Type toggle -->
				<div class="flex rounded-lg border border-border-input overflow-hidden">
					<button
						type="button"
						onclick={() => (newAllocationType = "shared")}
						class="flex-1 py-2 text-sm font-medium transition-colors {newAllocationType === 'shared'
							? 'bg-foreground text-background'
							: 'bg-background text-foreground-alt hover:bg-muted'}"
					>
						Partagée
					</button>
					<button
						type="button"
						onclick={() => (newAllocationType = "dedicated")}
						class="flex-1 py-2 text-sm font-medium transition-colors {newAllocationType === 'dedicated'
							? 'bg-foreground text-background'
							: 'bg-background text-foreground-alt hover:bg-muted'}"
					>
						Dédiée
					</button>
				</div>
				<!-- Partner select -->
				<div>
					<label for="ca-partner" class="block text-sm font-medium text-foreground mb-1">Partenaire *</label>
					<select id="ca-partner" bind:value={newAllocationPartnerId} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm">
						<option value="">Sélectionner un partenaire</option>
						{#each verifiedPartners as p (p.id)}
							<option value={p.user_id}>{p.user_id.slice(0, 8)}…</option>
						{/each}
					</select>
					{#if verifiedPartners.length === 0}
						<p class="mt-1 text-xs text-amber-600 dark:text-amber-400">Aucun partenaire vérifié disponible</p>
					{/if}
				</div>
				{#if newAllocationType === "dedicated"}
					<div>
						<label for="ca-start" class="block text-sm font-medium text-foreground mb-1">Date de début *</label>
						<input id="ca-start" type="date" bind:value={newAllocationStartDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
					<div>
						<label for="ca-end" class="block text-sm font-medium text-foreground mb-1">Date de fin</label>
						<input id="ca-end" type="date" bind:value={newAllocationEndDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
				{/if}
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (createAllocationOpen = false)} class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm">Annuler</button>
					<button onclick={submitCreateAllocation} disabled={loading["create-allocation"]} class="px-4 py-2 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm disabled:opacity-50">
						{loading["create-allocation"] ? "Création…" : "Créer"}
					</button>
				</div>
			</div>
		</Drawer>
	{:else}
		<Modal bind:isOpen={createAllocationOpen} title="Nouvelle allocation" description="Cabinet « {selectedRoom.name} »">
			<div class="mt-4 space-y-4">
				<div class="flex rounded-lg border border-border-input overflow-hidden">
					<button
						type="button"
						onclick={() => (newAllocationType = "shared")}
						class="flex-1 py-2 text-sm font-medium transition-colors {newAllocationType === 'shared'
							? 'bg-foreground text-background'
							: 'bg-background text-foreground-alt hover:bg-muted'}"
					>
						Partagée
					</button>
					<button
						type="button"
						onclick={() => (newAllocationType = "dedicated")}
						class="flex-1 py-2 text-sm font-medium transition-colors {newAllocationType === 'dedicated'
							? 'bg-foreground text-background'
							: 'bg-background text-foreground-alt hover:bg-muted'}"
					>
						Dédiée
					</button>
				</div>
				<div>
					<label for="ca-partner" class="block text-sm font-medium text-foreground mb-1">Partenaire *</label>
					<select id="ca-partner" bind:value={newAllocationPartnerId} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm">
						<option value="">Sélectionner un partenaire</option>
						{#each verifiedPartners as p (p.id)}
							<option value={p.user_id}>{p.user_id.slice(0, 8)}…</option>
						{/each}
					</select>
					{#if verifiedPartners.length === 0}
						<p class="mt-1 text-xs text-amber-600 dark:text-amber-400">Aucun partenaire vérifié disponible</p>
					{/if}
				</div>
				{#if newAllocationType === "dedicated"}
					<div>
						<label for="ca-start" class="block text-sm font-medium text-foreground mb-1">Date de début *</label>
						<input id="ca-start" type="date" bind:value={newAllocationStartDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
					<div>
						<label for="ca-end" class="block text-sm font-medium text-foreground mb-1">Date de fin</label>
						<input id="ca-end" type="date" bind:value={newAllocationEndDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
					</div>
				{/if}
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (createAllocationOpen = false)} class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium">Annuler</button>
					<button onclick={submitCreateAllocation} disabled={loading["create-allocation"]} class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium disabled:opacity-50">
						{loading["create-allocation"] ? "Création…" : "Créer"}
					</button>
				</div>
			</div>
		</Modal>
	{/if}
{/if}

<!-- Edit Period Modal -->
{#if editPeriodOpen && selectedAllocation}
	{#if isMobile}
		<Drawer bind:isOpen={editPeriodOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">Modifier la période</h2>
					<button type="button" onclick={() => (editPeriodOpen = false)} class="p-2 hover:bg-muted rounded-md transition-colors">
						<X class="text-foreground size-5" />
					</button>
				</div>
			</div>
			<div class="pt-4 space-y-4">
				<div>
					<label for="ep-start" class="block text-sm font-medium text-foreground mb-1">Date de début *</label>
					<input id="ep-start" type="date" bind:value={editPeriodStartDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="ep-end" class="block text-sm font-medium text-foreground mb-1">Date de fin</label>
					<input id="ep-end" type="date" bind:value={editPeriodEndDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (editPeriodOpen = false)} class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm">Annuler</button>
					<button onclick={submitEditPeriod} disabled={loading["edit-period"]} class="px-4 py-2 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium text-sm disabled:opacity-50">
						{loading["edit-period"] ? "Enregistrement…" : "Enregistrer"}
					</button>
				</div>
			</div>
		</Drawer>
	{:else}
		<Modal bind:isOpen={editPeriodOpen} title="Modifier la période">
			<div class="mt-4 space-y-4">
				<div>
					<label for="ep-start" class="block text-sm font-medium text-foreground mb-1">Date de début *</label>
					<input id="ep-start" type="date" bind:value={editPeriodStartDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div>
					<label for="ep-end" class="block text-sm font-medium text-foreground mb-1">Date de fin</label>
					<input id="ep-end" type="date" bind:value={editPeriodEndDate} class="w-full px-3 py-2 border border-border-input rounded-lg bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-foreground/20 text-sm" />
				</div>
				<div class="flex w-full justify-end gap-3 pt-2">
					<button type="button" onclick={() => (editPeriodOpen = false)} class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium">Annuler</button>
					<button onclick={submitEditPeriod} disabled={loading["edit-period"]} class="px-6 py-2.5 bg-foreground text-background rounded-lg hover:bg-foreground/90 transition-colors font-medium disabled:opacity-50">
						{loading["edit-period"] ? "Enregistrement…" : "Enregistrer"}
					</button>
				</div>
			</div>
		</Modal>
	{/if}
{/if}

<!-- Deactivate Confirmation Modal -->
{#if deactivateConfirmOpen && selectedAllocation}
	{#if isMobile}
		<Drawer bind:isOpen={deactivateConfirmOpen}>
			<div class="sticky top-0 bg-background pb-4 border-b border-border-card -mx-4 px-4 -mt-4 pt-4 z-10">
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xl font-semibold tracking-tight">Désactiver l'allocation</h2>
					<button type="button" onclick={() => (deactivateConfirmOpen = false)} class="p-2 hover:bg-muted rounded-md transition-colors">
						<X class="text-foreground size-5" />
					</button>
				</div>
				<p class="text-muted-foreground text-sm">
					Êtes-vous sûr de vouloir désactiver cette allocation ({selectedAllocation.allocation_type}) ?
				</p>
			</div>
			<div class="pt-4">
				<div class="flex w-full justify-end gap-3">
					<button type="button" onclick={() => (deactivateConfirmOpen = false)} class="px-4 py-2 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium text-sm">Annuler</button>
					<button onclick={submitDeactivate} disabled={loading["deactivate"]} class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium text-sm disabled:opacity-50">
						{loading["deactivate"] ? "Désactivation…" : "Désactiver"}
					</button>
				</div>
			</div>
		</Drawer>
	{:else}
		<Modal
			bind:isOpen={deactivateConfirmOpen}
			title="Désactiver l'allocation"
			description="Êtes-vous sûr de vouloir désactiver cette allocation ({selectedAllocation.allocation_type}) ?"
		>
			<div class="mt-6">
				<div class="flex w-full justify-end gap-3">
					<button type="button" onclick={() => (deactivateConfirmOpen = false)} class="px-6 py-2.5 border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors font-medium">Annuler</button>
					<button onclick={submitDeactivate} disabled={loading["deactivate"]} class="px-6 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium disabled:opacity-50">
						{loading["deactivate"] ? "Désactivation…" : "Désactiver"}
					</button>
				</div>
			</div>
		</Modal>
	{/if}
{/if}
