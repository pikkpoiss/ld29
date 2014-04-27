# ld29

Ludum Dare 29. Contributors:

  * Arne Roomann-Kurrik
  * Kalev Roomann-Kurrik
  * Wes Goodman

Theme: "Beneath the surface"

## Tasks

  * [x] Come up with idea.
  * [x] Main game loop.
  * [x] Level loading.
  * [x] Vertical level switching.
  * [x] Music loading.
  * [x] Player movement.
  * [x] Collision detection.
  * [x] Water tracking / gauge.
  * [x] Items.
  * [x] Pump item type and pump events.
  * [x] Menus.
  * [x] Load player sprite sheet.
  * [x] Triggers for running over items and transition tiles.
  * [x] Destructable item framework.
  * [x] Health bar.
  * [ ] Show picked up items.
  * [ ] End state.
  * [ ] Splash screen.
  * [ ] Level tiles for dry, partial, drowned levels.
  * [ ] Sound effects.
  * [ ] Snorkel.
  * [ ] Pickaxe and barrier tiles.
  * [ ] Pump tile.

## Brainstorming

  * Ground.
  * Change over time.
  * Digging holes, top down, holes have different depths,
    can only go adjacent levels.
  * Submarine.
  * Underground pathways.
  * More depth than initially expected.
  * Tremors the game.
  * Character that has to duck out of the way.
  * Revealing something about NPCs.
  * Tremors from monster's point of view.
    * NPCs on rocks, running to some goal.
    * Can go around rocks, eat NPCs if you touch them.
    * Need to wait for them to expose themselves.
    * Pathfinding with preferences to stay on rock.
  * Dungeon, go deeper.
  * Motivations for going beneath the surface:
    * Escape from surface
    * Escape from below

  * Well Idea:
    * Start at the top of a well.
    * Raining.
    * Climb down and levels strip away.
    * Need to find a way down.
    * Powerup items: Move faster, hold breath longer, etc.
    * Well starts to fill with water, so you need to be fast.
    * Can return to top of well to pump water out.
    * Need to get to bottom to get more loot.
    * Further down, better loot.
    * Levels fill with water slowly, eventually you
      hold your breath, take damage, etc.

## Setup

Complete the setup steps for the twodee lib.

Also:

	go get -u github.com/kurrik/tmxgo

Run:

	git submodule init
	git submodule update
