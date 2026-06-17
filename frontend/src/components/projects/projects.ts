import { LitElement, PropertyValues, TemplateResult, html, nothing } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ElementBase } from "../elementBase";
import { insertCustomElement, removeElement } from "../../utils";
import htmx from "htmx.org";

@customElement("camb-projects")
class ProjectsComponent extends ElementBase {
    @property({ attribute: "data-curr-pid" })
    accessor currentPid: string = "N/A";
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    private openNewProjectModal(_: Event) {
        insertCustomElement(
            '<new-project-modal id="new-project-modal"></new-project-modal>',
            document.body,
        );
    }
    protected firstUpdated(_changedProperties: PropertyValues): void {
        this.processFakeSlots(["projects-list"])
        htmx.process(this.renderRoot)
    }
    protected render(): TemplateResult<1> {
        return html`<div class="columns">
            <div class="column is-3">
                <nav class="panel">
                    <p class="panel-heading">Projects</p>
                    <div class="panel-block">
                        <button id="new-project" class="button is-link is-outlined is-fullwidth"
                        @click=${this.openNewProjectModal}>
                            New Project
                        </button>
                    </div>
                    <!-- // <div class=" panel-block">
                    // <p class="control has-icons-left">
                        // <input class="input" type="text" placeholder="Search" />
                        // <span class="icon is-left">
                            // <i class="fas fa-search" aria-hidden="true"></i>
                            // </span>
                        // </p>

                    //
            </div>
            // <p class="panel-tabs">
                // <a class="is-active">Last used</a>
                // <a>Last Modified</a>
                // </p> -->
                    <!-- for _, p := range projects {
            <a class="panel-block" hx-get={ projectLink(p.Id) } hx-target="#project-details">{
                p.Name }</a>
            } -->
                    <slot id="projects-list" name="projects-list"></slot>
                </nav>
            </div>
            <div id="project-details" class="column" hx-get=${"/components/project/" + this.currentPid} hx-trigger="load">Select or create a new project</div>
        </div> `;
    }
}

@customElement("camb-project")
class ProjectComponent extends ElementBase {
    @property({ attribute: "data-pid" })
    accessor projectId: string = "N/A";
    @property({ attribute: "data-name" })
    accessor projectName: string = "N/A";
    @property({ attribute: "data-desc" })
    accessor projectDesc: string = "N/A";
    @property({ attribute: "data-created" })
    accessor projectCreated: string = "N/A";
    @property({ attribute: "data-updated" })
    accessor projectUpdated: string = "N/A";
    protected firstUpdated(_changedProperties: PropertyValues): void {
        this.processFakeSlots(["version-stats"]);
    }
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    protected render(): TemplateResult<1> {
        return html`<div id="root" class="block">
                <h1 id="title" class="title">${this.projectName}</h1>
                <h2 class="subtitle">Id: ${this.projectId}</h2>
                <p>${this.projectDesc}</p>
                <p>Created At: ${this.projectCreated}</p>
                <p>Last Update: ${this.projectUpdated}</p>
                <h2 class="is-size-3">Versions:</h2>
                <nav class="level">
                    <slot name="version-stats"></slot>
                </nav>
                <form
                    id="edit-project"
                    hx-target="#project-details"
                    hx-on::before-request="hideElement('newproject')"
                >
                    <label class="label">Recipe</label>
                    <div class="select is-primary">
                        <select name="recipe"></select>
                    </div>
                    <div class="field is-grouped">
                        <input type="submit" class="button is-link" value="Update" />
                    </div>
                </form>
            </div>
            <div class="block">
                <nav>
                    <button
                        id="set-active"
                        class="button is-link"
                        hx-swap="outerHTML"
                        hx-target="global #main-body"
                    >
                        Set as active
                    </button>
                    <button id="btn-delete" class="button is-link" hx-target="global #main-body">
                        Delete
                    </button>
                </nav>
            </div> `;
    }
}

@customElement("new-project-modal")
class NewProjectModal extends LitElement {
    protected firstUpdated(_changedProperties: PropertyValues): void {
        htmx.process(this.renderRoot)
    }
    protected createRenderRoot(): HTMLElement | DocumentFragment {
        return this;
    }
    protected closeMe(_: Event) {
        removeElement("new-project-modal", document.body);
    }
    protected render(): TemplateResult {
        return html`<div id="newproject" class="modal is-active">
            <div class="modal-background"></div>
            <div class="modal-content">
                <div class="box">
                    <h3 class="is-size-4">Create new project</h3>
                    <form
                        id="new-project-submit"
                        hx-post="/components/project"
                        hx-swap="outerHTML"
                        hx-target="#main-body"
                        @htmx:afterRequest=${this.closeMe}>
                        <div class="field">
                            <label class="label">Project Name</label>
                            <input
                                class="input"
                                name="name"
                                type="text"
                                placeholder="Input project name"
                            />
                        </div>
                        <div class="field">
                            <input type="submit" class="button is-link" value="Create" />
                        </div>
                    </form>
                </div>
            </div>
            <button
                id="close-new-project-btn"
                class="modal-close is-large"
                aria-label="close"
            ></button>
        </div> `;
    }
}
