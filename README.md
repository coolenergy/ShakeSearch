# ShakeSearch

Welcome to ShakeSearch, a simple web app for searching the complete works of William Shakespeare. The app provides an easy-to-use interface to search for any text string within Shakespeare's writings.

## Check the live version here:
https://rocky-brook-64249.herokuapp.com/

Features and Improvements

The following improvements have been made to enhance the user experience and the app's functionality:

- Case-insensitive search: The search is now case-insensitive, allowing users to find matches regardless of the capitalization used in the query.

- Highlight search terms: The search terms in the results are highlighted for better readability and easy identification of the matched text.

- Whitespace normalization: Extra whitespace and line breaks in the search results have been normalized for improved readability.

- Handling multi-word and partial matches: The search algorithm has been improved to handle multi-word queries and partial matches, providing more relevant results to users.

- Pagination: Pagination has been implemented to manage large numbers of search results, allowing users to navigate through the results more easily.

- Improved search speed: The search algorithm has been optimized to increase the speed and efficiency of the search, providing faster results to users.
Usage

To search for a text string, simply enter the desired text in the input field and click the "Search" button. The app will return a list of results containing the search terms. Use the pagination controls to navigate through the search results if there are multiple pages.

### API
The search functionality is exposed through an API endpoint:
GET /search?q=<query>&page=<page_number>&perPage=<results_per_page>

### Parameters:
q: The search query (required)
page: The page number (optional, default: 1)
perPage: The number of results per page (optional, default: 10)
The API returns a JSON object containing the search results, current page, total pages, and per-page count.

### Running Locally

To run the ShakeSearch app locally, follow these steps:

Clone the repository to your local machine.
Install Go if you haven't already.
Run go run main.go in the terminal to start the server.
Open your browser and navigate to http://localhost:3001.

## Enjoy searching through the works of Shakespeare!