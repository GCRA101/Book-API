<script>
export default {
  data: function() {
    return {
      // Declare one reactive array of books
      books: [],
   }
  },
  methods: {
    // Declare a function that sends an HTTP GET request via Axios
    async loadBooksList() {
      try {
        let response = await this.$axios.get("/books");
        console.log("Response:", response);
        this.books=(response.data).data;
      } catch (err) {
        alert("Error: " + err);
      }
    }
  }
}
</script>

<template>
  <!-- When this button is clicked, the list of books is downloaded -->
  <div> 
  <button @click="loadBooksList">
      Load Books
  </button>
  </div>
  <div>
    <ul>
        <!-- Loop the <li> tag for each book -->
        <li v-for="b in books">Title: {{b.title}} - Author:{{b.author}}</li>
    </ul>
   </div>
</template>
