import { LitElement, PropertyValues, html, nothing } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ElementBase } from "../elementBase";
import htmx from "htmx.org";

@customElement("camb-specs")
class SpecsComponent extends ElementBase {
    protected render() {
        this.copyGlobalStyles();
        return html`<div class="box">
                <form id="new-version-form">
                    <div class="container is-fluid">
                        <b>Create a new version / epic:</b>
                        <label class="cbs_hform_item" for="version_name">Name</label>
                        <input
                            class="cbs_hform_item"
                            type="text"
                            id="version_name"
                            name="version_name"
                        />
                        <input class="cbs_hform_item" type="submit" value="create" />
                    </div>
                </form>
            </div>
            <div class="box cbs_vfull" id="versionlist-container">
                <slot name="version-list"></slot>
            </div>`;
    }
}

@customElement("version-item")
class VersionItem extends ElementBase {
    @property({ attribute: "data-name" })
    accessor versionName: string = "N/A";
    @property({ attribute: "data-desc" })
    accessor versionDesc: string = "N/A";
    @property({ attribute: "data-status" })
    accessor versionStatus: string = "N/A";

    protected firstUpdated(_changedProperties: PropertyValues): void {
        this.processFakeSlots(["vi-story-status"]);
    }
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    protected render() {
        return html`<div name="vi-root" class="box container is-fluid">
            <button
                name="vi-delete"
                class="delete is-pulled-right"
                hx-target="closest div"
                hx-confirm="This will permanently delete the version. Are you sure?"
                hx-swap="outerHTML"
            ></button>
            <div class="level">
                <div class="level-item has-text-centered">
                    <div>
                        <p class="heading">Name</p>
                        <p><b name="vi-name">${this.versionName}</b></p>
                    </div>
                </div>
                <div class="level-item has-text-centered">
                    <div>
                        <p class="heading">Description</p>
                        <p><b name="vi-desc">${this.versionDesc}</b></p>
                    </div>
                </div>
                <div class="level-item has-text-centered">
                    <div>
                        <p class="heading">Status</p>
                        <p><b name="vi-status">${this.versionStatus}</b></p>
                    </div>
                </div>
                <div class="level-item has-text-centered">
                    <div>
                        <button name="vi-edit" class="button">Edit</button>
                    </div>
                </div>
                <slot name="vi-story-status"></slot>
            </div>
        </div>`;
    }
}

@customElement("edit-version-modal")
class EditVersionModal extends LitElement {
    protected firstUpdated(_changedProperties: PropertyValues): void {
        htmx.process(this.renderRoot)
    }
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    protected render() {
        return html`<div id="editversion" class="modal is-active">
            <div function="close-edit-version" class="modal-background"></div>
            <div class="modal-content">
                <div class="box">
                    <h3 class="is-size-4">Edit version</h3>
                </div>
                <button
                    function="close-edit-version"
                    class="modal-close is-large"
                    aria-label="close"
                ></button>
            </div>
        </div> `;
    }
}
