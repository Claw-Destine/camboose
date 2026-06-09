import { FieldMapping, ShadowTemplElement, TemplElement } from '../../elementBase';
import projectsTemplate from './camb-projects.html';
import projectTemplate from './camb-project.html';
import newProjectModalTemplate from './new-project-modal.html'
import { registerElementWithTemplate } from '../../elementBase';
import htmx from 'htmx.org';
import { insertCustomElement, removeElement } from '../../utils';

export class ProjectsComponent extends ShadowTemplElement {
    constructor() {
        super("camb-projects", true);
        this.wireButtons();
    }

    wireButtons() {
        const npb = this.shadowRoot.querySelector('button[id="new-project"]') as HTMLButtonElement | null;
        npb.addEventListener('click', _ => { insertCustomElement('<new-project-modal id="new-project-modal"></new-project-modal>', document.body) })
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

export class NewProjectModal extends TemplElement {
    constructor() {
        super("new-project-modal")
        const closeBtn = this.querySelector('button[id="close-new-project-btn"]');
        closeBtn.addEventListener("click", _ => { removeElement("new-project-modal", document.body) })
        const npForm = this.querySelector('form[id="new-project-submit"]');
        closeBtn.addEventListener('htmx:before-request', _ => { removeElement("new-project-modal", document.body) })
    }
}

export function registerNewProjectModal() {
    registerElementWithTemplate("new-project-modal", NewProjectModal, newProjectModalTemplate)
}