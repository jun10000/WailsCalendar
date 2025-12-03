import { GetCalendar, SetCalendar, OffsetCalendar } from './wailsjs/go/main/App';

function updateCalendar(source) {
    const labelElement = document.getElementById("calendar-label");
    const dayElements = document.getElementById("calendar-days");
    const current = source.Current.Date;

    labelElement.textContent = current.Year.toString().concat("/", current.Month.toString());
    dayElements.innerHTML = "";
    for (let d of source.Days) {
        const e = document.createElement("div");
        e.className = d.ClassName;
        e.textContent = d.Label;
        e.addEventListener("click", (ev) => {
            SetCalendar(d.Date.Year, d.Date.Month, d.Date.Day).then((c) => updateCalendar(c)).catch((r) => console.error(r));
        });
        dayElements.appendChild(e);
    }
}

window.addEventListener("DOMContentLoaded", (e) => {
    const prevButtonElement = document.getElementById("calendar-prev-button");
    const nextButtonElement = document.getElementById("calendar-next-button");
    const todayButtonElement = document.getElementById("calendar-today-button");

    GetCalendar().then((c) => updateCalendar(c)).catch((r) => console.error(r));

    prevButtonElement.addEventListener("click", (e) => {
        OffsetCalendar(0, -1, 0, false).then((c) => updateCalendar(c)).catch((r) => console.error(r));
    });
    nextButtonElement.addEventListener("click", (e) => {
        OffsetCalendar(0, 1, 0, false).then((c) => updateCalendar(c)).catch((r) => console.error(r));
    });
    todayButtonElement.addEventListener("click", (e) => {
        const today = new Date();
        const [y, m, d] = [today.getFullYear(), today.getMonth() + 1, today.getDate()];
        SetCalendar(y, m, d).then((c) => updateCalendar(c)).catch((r) => console.error(r));
    });
});
