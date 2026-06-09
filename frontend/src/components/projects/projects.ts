import { FieldMapping, ShadowTemplElement } from '../../elementBase';
import projectsTemplate from './camb-projects.html';
import projectTemplate from './camb-project.html';
import { registerElementWithTemplate } from '../../elementBase';
import htmx from 'htmx.org';

export class ProjectsComponent extends ShadowTemplElement {
    constructor() {
        super("camb-projects", true);
    }


}

export function registerProjectsComponent() {
    registerElementWithTemplate("camb-projects", ProjectsComponent, projectsTemplate);
}
const projectMappings: FieldMapping[] = [
    { source: "data-name", targetSelector: "#title" },
    { source: "data-pid", targetSelector: "#pid" },
    { source: "data-created", targetSelector: "#created" },
    { source: "data-updated", targetSelector: "#updated" },
    { source: "data-recipe", targetSelector: "#recipies", isList: true }
];
export class ProjectComponent extends ShadowTemplElement {

    constructor() {
        super("camb-project", true, projectMappings);
        this.populateRecipiesList();
    }

    populateRecipiesList() {
        const recipeSelect = this.shadowRoot?.querySelector('select[name="recipe"]') as HTMLSelectElement | null;
        if (!recipeSelect) {
            return;
        }

        recipeSelect.innerHTML = "";
        const attrName = "data-curr-recipe";
        let currRecipe = ""
        if (this.hasAttribute(attrName)) {
            currRecipe = this.getAttribute(attrName);
        }
        const option = document.createElement("option");
        option.value = currRecipe;
        option.textContent = currRecipe;
        recipeSelect.appendChild(option);


        for (let idx = 0; ; idx++) {
            const attrName = `data-recipe-${idx}`;
            if (!this.hasAttribute(attrName)) {
                break;
            }

            const value = this.getAttribute(attrName) ?? "";
            if (value === currRecipe) {
                break;
            }
            const option = document.createElement("option");
            option.value = value;
            option.textContent = value;
            recipeSelect.appendChild(option);
        }
    }
}

export function registerProjectComponent() {
    registerElementWithTemplate("camb-project", ProjectComponent, projectTemplate);
}