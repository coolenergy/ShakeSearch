function highlightSearchTerms(text, query) {
  const queryTerms = query.split(/\s+/);
  const escapedTerms = queryTerms.map(term => term.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&'));
  const regex = new RegExp(`(${escapedTerms.join('|')})`, 'gi');
  const cleanedText = text.replace(/[_\[\]]/g, "").replace(/\s\s+/g, " ");
  return cleanedText.replace(regex, '<mark>$1</mark>');
}

const Controller = {
  search: async (ev, page = 1) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = await fetch(`/search?q=${data.query}&page=${data.page}&perPage=${data.perPage}`);
    const results = await response.json();
    Controller.updateTable(results);
},

  updateTable: (results) => {
    const query = document.getElementById("query").value;
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results.results) {
      const formattedResult = result.replace(/(\r\n|\n|\r)/gm, "<br>").replace(/\s\s+/g, " ");
      const highlightedResult = highlightSearchTerms(formattedResult, query);
      rows.push(`<tr><td>${highlightedResult}</td></tr>`);
    }
    table.innerHTML = rows.join("");

    // Update pagination
    const pagination = document.getElementById("pagination");
    let pages = '';
    for (let i = 1; i <= results.totalPages; i++) {
      pages += `<button onclick="Controller.search(event, ${i})">${i}</button>`;
    }
    pagination.innerHTML = pages;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", (event) => Controller.search(event));