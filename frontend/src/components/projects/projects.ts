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
const projectMappings = [
    { source: "data-name", targetSelector: "#title", targetAttribute: null },
    { source: "data-pid", targetSelector: "#pid", targetAttribute: null },
    { source: "data-created", targetSelector: "#created", targetAttribute: null },
    { source: "data-updated", targetSelector: "#updated", targetAttribute: null },
];
export class ProjectComponent extends ShadowTemplElement {

    constructor() {
        super("camb-project", true, projectMappings);

    }
}

export function registerProjectComponent() {
    registerElementWithTemplate("camb-project", ProjectComponent, projectTemplate);
}