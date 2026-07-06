# PRD-0013: Environment Management UI

## Problem Statement

I cannot create, edit, or delete environments inside the graphical user interface. The environment list loader panics on execution, and there is no management panel on screen, forcing me to manually edit and manage YAML files in my file manager to configure environments.

## Solution

Implement an Environment Management interface in the GUI. Provide a modal **Environment Manager** dialog accessed via the top menu bar or a configuration button next to the environment selector dropdown. The manager will feature a split-pane layout: a list of environments on the left with Add/Delete buttons, and an editor pane on the right containing a name input and a Dynamic Row List to configure environment variables.

## User Stories

1. As a user, I want to manage my environment profiles in a dedicated pop-up modal, so that variable editing doesn't clutter my request collection workspace.
2. As a user, I want to open the Environment Manager by selecting `Environment > Manage Environments...` in the top menu or clicking a gear button next to the active environment dropdown in the header.
3. As a user, I want to create new environments and delete existing ones within the manager dialog, so that I can easily set up configurations for new projects or environments.
4. As a user, I want to rename the selected environment, so that its name matches its purpose (e.g. changing "New Env" to "Staging").
5. As a user, I want to edit the variables (key-value pairs) of my selected environment using a Dynamic Row List, so that I can add or update parameters visually.
6. As a user, I want my changes in the Environment Manager to be persisted back to disk as separate YAML files in my environment storage directory on save, so that my environments are permanently saved.

## Implementation Decisions

* **Manager Layout:**
  * Implement the Environment Manager using a modal window wrapper (`fyne.Window` or `widget.NewModal`).
  * Left pane: A `widget.NewList` displaying the names of all loaded environments. Below the list, place two buttons: **Add** (plus icon) and **Delete** (trash icon).
  * Right pane: A vertical layout containing:
    * Environment Name: `widget.NewEntry` (bound to the selected environment's name).
    * Variables Editor: A Dynamic Row List (as defined in PRD-0012) bound to the selected environment's variables map.
* **Saving Changes:**
  * Provide **Cancel** and **Save** action buttons at the bottom of the modal.
  * Clicking **Save** runs validation checks, writes the updated environment structs back to disk (overwriting `<env_name>.yaml` files), reloads the active environment list in the main header dropdown, and closes the modal.
  * Clicking **Cancel** discards in-memory modifications and closes the modal.

## Testing Decisions

* **Testing Philosophy:**
  * Test the Environment Manager panel by using Fyne's UI test runner. Instantiate the dialog programmatically, simulate button clicks (Add, Delete), enter characters into the inputs, and assert that the underlying model state is modified correctly and save triggers the persistence write functions.
* **Modules to Test:**
  * `internal/ui/environments` (the manager UI controller).
* **Prior Art:**
  * None.

## Out of Scope

* Importing variables from `.env` files or JSON files inside this manager (generic importing is covered in the OpenAPI Import PRD).
* Specifying variable types (all variables are stored as text and cast automatically by the substitution engine).
