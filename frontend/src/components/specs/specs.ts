import { LitElement, PropertyValues, html, nothing } from "lit";
import { customElement, property, state } from "lit/decorators.js";
import { ElementBase } from "../elementBase";
import htmx from "htmx.org";
import { insertCustomElement, removeElement, showNotification } from "../../utils";

@customElement("camb-specs")
class SpecsComponent extends ElementBase {
    @property({ attribute: "data-pid" })
    accessor pid: string;
    protected updated(_changedProperties: PropertyValues): void {
        this.processFakeSlots(["version-list"]);
        htmx.process(this.renderRoot);
    }
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    protected render() {
        this.copyGlobalStyles();
        return html`<div class="box">
                <form
                    id="new-version-form"
                    hx-post=${"/components/version?currentProject=" + this.pid}
                    hx-target="#versionlist-container"
                    @htmx:responseError=${showNotification}
                >
                    <div class="container is-fluid">
                        <b>Create a new version / epic:</b>
                        <label class="cbs_hform_item" for="version_name">Name: </label>
                        <input
                            class="cbs_hform_item"
                            type="text"
                            id="version_name"
                            name="version_name"
                        />
                        <input hidden type="text" name="pid" value=${this.pid} />
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
    @property({ attribute: "data-id" })
    accessor vid: string;
    @property({ attribute: "data-name" })
    accessor versionName: string = "N/A";
    @property({ attribute: "data-desc" })
    accessor versionDesc: string = "N/A";
    @property({ attribute: "data-status" })
    accessor versionStatus: string = "N/A";
    @property({ attribute: "data-stats", type: Object })
    accessor stats: Object
    @state()
    accessor verboseVersion: boolean = false

    protected updated(_changedProperties: PropertyValues): void {
        this.processFakeSlots(["vi-story-status"]);
    }
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    private openEditVersuibModal(_: Event) {
        insertCustomElement(
            '<edit-version-modal id="edit-version-modal"></edit-version-modal>',
            document.body,
        );
    }
    private expandItem = () => {
        this.verboseVersion = true
    }
    private retractItem = () => {
        this.verboseVersion = false
    }
    private renderBrief() {
        const statsTemplates = [];
        for (const k in this.stats) {
            statsTemplates.push(html`<div class="level-item has-text-centered">
                    <div>
                        <p class="heading">${k}</p>
                        <p><b name="vi-desc">${this.stats[k]}</b></p>
                    </div>
                </div>`)
        }

        return html`<div
            name="vi-root"
            class="box container is-fluid"
            @click=${this.expandItem}
        >
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
                <!-- <div class="level-item has-text-centered">
                    <div>
                        <button name="vi-edit" class="button" @click=${this.openEditVersuibModal}>
                            Edit
                        </button>
                    </div>
                </div> -->
                ${statsTemplates}
                <slot name="vi-story-status"></slot>
            </div>
        </div>`;
    }
    protected renderVerbose() {
        return html`<div class="box columns">
            <div class="column is-one-fifth">
                <div class="level">
                    <div class="level-item has-text-centered">
                        <div>
                            <button name="vi-edit" class="button"
                            @click=${this.retractItem}>Close</button>
                        </div>
                    </div>
                    <div class="level-item has-text-centered">
                        <div>
                            <button name="vi-edit" class="button">Edit</button>
                        </div>
                    </div>
                    <div class="level-item has-text-centered">
                        <div>
                            <button name="vi-edit" class="button">Delete</button>
                        </div>
                    </div>

                </div>
                <div class="field">
                    <label class="label">Name</label>
                    <div class="control">
                        <input class="input" type="text" placeholder=${this.versionName} />
                    </div>
                </div>
                <div class="field">
                    <label class="label">Description</label>
                    <div class="control">
                        <input class="input" type="text" placeholder=${this.versionDesc} />
                    </div>
                </div>
            </div>
            <div class="column">Stories</div>
        </div>`;
    }
    protected render() {
        if (this.verboseVersion) {
            return html`${this.renderVerbose()}`;
        } else {
            return html`${this.renderBrief()}`;
        }
    }
}

@customElement("edit-version-modal")
class EditVersionModal extends LitElement {
    protected firstUpdated(_changedProperties: PropertyValues): void {
        htmx.process(this.renderRoot);
    }
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    protected closeMe(_: Event) {
        removeElement("#edit-version-modal", document.body);
    }
    protected render() {
        return html`<div id="editversion" class="modal is-active">
            <div class="modal-background" @click=${this.closeMe}></div>
            <div class="modal-content">
                <div class="box">
                    <h3 class="is-size-4">Edit version</h3>
                </div>
                <button
                    class="modal-close is-large"
                    aria-label="close"
                    @click=${this.closeMe}
                ></button>
            </div>
        </div> `;
    }
}
