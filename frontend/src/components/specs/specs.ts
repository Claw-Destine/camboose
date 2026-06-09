import { registerElementWithTemplate, ShadowTemplElement } from "../../elementBase";
import specsTemplate from "./camb-specs.html"

export class SpecsComponent extends ShadowTemplElement {
    constructor() {
        super("camb-specs", true);
        // this.wireButtons();
    }
}

export function registerSpecsComponent() {
    registerElementWithTemplate("camb-specs", SpecsComponent, specsTemplate);
}