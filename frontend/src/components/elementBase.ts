import { LitElement } from "lit";

export class ElementBase extends LitElement {
    protected copyGlobalStyles() {
        const globalStyles = document.querySelectorAll("style"); // or any identifier
        globalStyles.forEach((style) => {
            this.renderRoot?.appendChild(style.cloneNode(true));
        });
    }

    protected processFakeSlots(fakeSlots: string[]) {
        for (const sn of fakeSlots) {
            const slot = this.renderRoot?.querySelector('slot[name="' + sn + '"]')?.parentNode;
            const slottables = this.renderRoot?.querySelectorAll("[slot=" + sn + "]");
            for (const sl of slottables) {
                slot.appendChild(sl);
            }
        }
    }
}
