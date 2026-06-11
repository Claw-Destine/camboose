import htmx from "htmx.org";
import { FieldMapping, registerElementWithTemplate, ShadowTemplElement, TemplElement } from "../../elementBase";
import specsTemplate from "./camb-specs.html"
import versionItemTemplate from "./version-item.html"
import editVersionModalTemplate from "./edit-version-modal.html"
import { insertCustomElement, removeElement } from "../../utils";



export class SpecsComponent extends ShadowTemplElement {
    constructor() {
        super("camb-specs", true);
        this.setupActions()
    }

    setupActions() {
        const pid = this.getAttribute("data-project")
        const form = this.shadowRoot.querySelector('form[id="new-version-form"]')
        const targetElement = this.shadowRoot.querySelector('div[id="versionlist-container"]')

        function createVersion(event: SubmitEvent) {
            event.preventDefault();

            const data = new FormData(event.target as HTMLFormElement);
            const values: Record<string, FormDataEntryValue> = {};
            for (const [key, value] of data.entries()) {
                values[key] = value;
            }
            values["pid"] = pid
            htmx.ajax("POST", "/components/version/?currentProject=" + pid,
                {
                    "target": targetElement,
                    "values": values
                }
            )
        }
        form.addEventListener("submit", createVersion)
    }
}

export function registerSpecsComponent() {
    registerElementWithTemplate("camb-specs", SpecsComponent, specsTemplate);
}
const versionMapping: FieldMapping[] = [{
    source: "data-name",
    targetSelector: 'b[name="vi-name"]'
}, {
    source: "data-desc",
    targetSelector: 'b[name="vi-desc"]'
}, {
    source: "data-status",
    targetSelector: 'b[name="vi-status"]'
}
]
export class VersionItem extends TemplElement {
    constructor() {
        super("version-item", versionMapping, ["vi-story-status"])
        const vid = this.getAttribute("data-id")
        const delBtn = this.querySelector('button[name="vi-delete"]')
        delBtn.setAttribute("hx-delete", "/components/version/" + vid)
        const npb = this.querySelector('button[name="vi-edit"]') as HTMLButtonElement | null;
        npb.addEventListener('click', _ => {
            console.log("edit version");
            insertCustomElement('<edit-version-modal id="edit-version-modal"></edit-version-modal>',
                document.body);
        })
    }
}

export function registerVersionItem() {
    registerElementWithTemplate("version-item", VersionItem, versionItemTemplate)
}
const EDIT_VERSION_MODAL = "edit-version-modal"
export class EditVersionModal extends TemplElement {
    constructor() {
        super(EDIT_VERSION_MODAL)

        const closeBtn = this.querySelectorAll('[function="close-edit-version"]');
        closeBtn.forEach(item => item.addEventListener("click", _ => { removeElement(EDIT_VERSION_MODAL, document.body) }))
        htmx.process(this)
    }
}

export function registerEditVersionModal() {
    registerElementWithTemplate(EDIT_VERSION_MODAL, EditVersionModal, editVersionModalTemplate)
}