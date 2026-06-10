import { FieldMapping, registerElementWithTemplate, ShadowTemplElement, TemplElement } from "../../elementBase";
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
const specMappings: FieldMapping[] = [{
    source: "data-name",
    targetSelector: 'b[name="vi-name"]'
}, {
    source: "data-desc",
    targetSelector: 'b[name="vi-desc"]'
}, {
    source: "data-status",
    targetSelector: 'b[name="vi-status"]'
}]
export class VersionItem extends TemplElement {
    constructor() {
        super("version-item", specMappings, ["vi-story-status"])
    }
}

export function registerVersionItem() {
    registerElementWithTemplate("version-item", VersionItem, versionItemTemplate)
}