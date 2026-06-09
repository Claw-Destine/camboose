import { registerElementWithTemplate, ShadowTemplElement, TemplElement } from "../../elementBase";
import specsTemplate from "./camb-specs.html"
import versionItemTemplate from "./version-item.html"

export class SpecsComponent extends ShadowTemplElement {
    constructor() {
        super("camb-specs", true);
        // this.wireButtons();
    }
}

export function registerSpecsComponent() {
    registerElementWithTemplate("camb-specs", SpecsComponent, specsTemplate);
}

export class VersionItem extends TemplElement {
    constructor() {
        super("version-item")
    }
}

export function registerVersionItem() {
    registerElementWithTemplate("version-item", VersionItem, versionItemTemplate)
}