import { LitElement, html, nothing } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import "bulma/css/bulma.min.css";
import logo from "../../assets/camboose_logo.png"
import { getCookie, toggleClass } from "../../utils"
import htmx from 'htmx.org';
import { ifDefined } from 'lit/directives/if-defined.js';


@customElement('camb-body')
class MyElement extends LitElement {
    @property()
    accessor pid: string;
    @property()
    accessor view: string;

    private setActiveMenu(event: Event) {
        const menuItems = this.renderRoot?.querySelectorAll<HTMLAnchorElement>('#menu-list a');
        if (!menuItems) {
            return;
        }
        menuItems.forEach(item => item.classList.remove('is-active'));
        (event.target as Element).classList.add('is-active');
    }

    private setupQuickProjectMenu() {
        const button = this.renderRoot?.querySelector<HTMLElement>('#quick-project-button');
        if (!button || !this.renderRoot) {
            return;
        }

        button.addEventListener('click', _ => toggleClass(this.renderRoot as ShadowRoot, "quick-project-select"));
        const caption = button.querySelector<HTMLElement>("#quick-project-caption");
        if (!caption) {
            return;
        }

        function setCaption(item: Node, root: ParentNode) {
            if (item instanceof Element && item.classList.contains("is-active")) {
                caption.textContent = item.textContent;
            }

        }

        const slot = this.renderRoot.querySelector<HTMLSlotElement>('#quick-project-slot');
        if (!slot) {
            return;
        }

        const projectLinks = slot.assignedNodes();
        projectLinks.forEach(item => setCaption(item, this.renderRoot));
    }

    protected updated(changeProperties) {
        htmx.process(this.renderRoot)
        this.setupQuickProjectMenu()
    }

    protected render() {
        const globalStyles = document.querySelectorAll('style'); // or any identifier
        globalStyles.forEach(style => {
            this.renderRoot?.appendChild(style.cloneNode(true));
        });

        let view = this.view
        if (!view) {
            view = getCookie("view")
        }
        let pid = this.pid
        if (!pid) {
            pid = getCookie("project")
        }
        const suffix = pid ? "?currentProject=" + pid : "";
        const currViewPath = "/components/" + view + suffix;
        const specsPath = "/components/specs" + suffix;
        const tasksPath = "/components/specs" + suffix;
        const projectsPath = "/components/specs" + suffix;
        const recipiesPath = "/components/specs" + suffix;

        return html`
<div class="columns">
    <aside class="menu column is-narrow">
        <figure class="image is-128x128">
            <img src=${logo} alt="logo" />
        </figure>
        <h1 class="title">camboose</h1>
        <p class="menu-label">Current Project</p>
        <div id="quick-project-select" class="dropdown">
            <div id="quick-project-button" class="dropdown-trigger">
                <button id="quick-project-caption" class="button" aria-haspopup="true" aria-controls="dropdown-menu">
                    Select project
                </button>
            </div>
            <div class="dropdown-menu" id="quick-project-menu" role="menu">
                <div class="dropdown-content">
                    <slot id="quick-project-slot" name="project-dropdown">Create new project</slot>
                </div>
            </div>
        </div>

        <ul class="menu-list" id="menu-list">
            <p class="menu-label">Plan</p>
            <li><a id="specs-menu-item" hx-trigger="click" hx-target="#main-container"
            @click=${this.setActiveMenu} class="${ifDefined(view === 'specs' ? 'is-active' : "")}" 
            hx-get=${specsPath}>Specs</a></li>
            <p class="menu-label">Deliver</p>
            <li><a id="tasks-menu-item" hx-trigger="click" hx-target="#main-container"
            @click=${this.setActiveMenu} class="${ifDefined(view === 'tasks' ? 'is-active' : "")}" 
             hx-get=${tasksPath}>Tasks</a></li>
            <p class="menu-label">Manage</p>
            <li><a id="projects-menu-item" hx-trigger="click" hx-target="#main-container"
            @click=${this.setActiveMenu} class="${ifDefined(view === 'projects' ? 'is-active' : "")}"  
            hx-get=${projectsPath}>Projects</a>
            </li>
            <li><a id="recipies-menu-item" hx-trigger="click" hx-target="#main-container"
            @click=${this.setActiveMenu} class="${ifDefined(view === 'recipies' ? 'is-active' : "")}" 
             hx-get=${recipiesPath}>Recipies</a>
            </li>
        </ul>
    </aside>
    <div id="main-container" class="column" hx-trigger="load" hx-get=${currViewPath}>
    </div>
</div>
    `;
    }
}