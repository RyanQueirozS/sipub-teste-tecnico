function parseLatLong(latLongStr) {
  // Split the input string by the comma
  const [lat, long] = latLongStr.split(",");

  // Convert the values to numbers
  const latitude = parseFloat(lat.trim());
  const longitude = parseFloat(long.trim());

  return { latitude, longitude };
}

document
  .getElementById("customerForm")
  .addEventListener("submit", async function (e) {
    e.preventDefault(); // Prevent the form from submitting normally

    const userData = {
      IsActive: true,
      IsDeleted: false,
      Email: document.getElementById("email").value,
      Cpf: document.getElementById("cpf").value,
      Name: document.getElementById("nome").value,
    };

    // Make sure the parseLatLong function is defined before it's called
    const { latitude, longitude } = parseLatLong(
      document.getElementById("geolocalizacao").value,
    );

    const addressData = {
      IsActive: true,
      IsDeleted: false,
      Street: document.getElementById("logradouro").value,
      Number: document.getElementById("numero").value,
      Neighborhood: document.getElementById("bairro").value,
      Complement: document.getElementById("complemento").value,
      City: document.getElementById("cidade").value,
      State: document.getElementById("estado").value,
      Country: document.getElementById("pais").value,
      Latitude: latitude,
      Longitude: longitude,
    };

    const productData = {
      IsActive: true,
      IsDeleted: false,
      WeightGrams: Number(document.getElementById("peso").value),
      Price: Number(document.getElementById("productPrice").value),
      Name: document.getElementById("productName").value,
    };
    console.log(productData);

    // USER
    try {
      // API endpoint (replace with your own API URL)
      const apiUrl = "http://localhost:8080/u";

      // Send the data to the API using POST method
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(userData), // Convert the form data to JSON
      });

      // Check if the response is OK
      if (!response.ok) {
        throw new Error("Failed to send data");
      }

      // Display a success message
      document.getElementById("responseMessage").innerHTML =
        `<p class="text-success">Dados enviados com sucesso!</p>`;
      document.getElementById("customerForm").reset(); // Reset the form after submission
    } catch (error) {
      // Display an error message
      document.getElementById("responseMessage").innerHTML =
        `<p class="text-danger">Erro: ${error.message}</p>`;
    }

    // Address
    try {
      // API endpoint (replace with your own API URL)
      const apiUrl = "http://localhost:8080/addresses";

      // Send the data to the API using POST method
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(addressData), // Convert the form data to JSON
      });

      // Check if the response is OK
      if (!response.ok) {
        throw new Error("Failed to send data");
      }

      // Display a success message
      document.getElementById("responseMessage").innerHTML =
        `<p class="text-success">Dados enviados com sucesso!</p>`;
      document.getElementById("customerForm").reset(); // Reset the form after submission
    } catch (error) {
      // Display an error message
      document.getElementById("responseMessage").innerHTML =
        `<p class="text-danger">Erro: ${error.message}</p>`;
    }

    // Product

    try {
      // API endpoint (replace with your own API URL)
      const apiUrl = "http://localhost:8080/products";

      // Send the data to the API using POST method
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(productData), // Convert the form data to JSON
      });

      // Check if the response is OK
      if (!response.ok) {
        throw new Error("Failed to send data");
      }

      // Display a success message
      document.getElementById("responseMessage").innerHTML =
        `<p class="text-success">Dados enviados com sucesso!</p>`;
      document.getElementById("customerForm").reset(); // Reset the form after submission
    } catch (error) {
      // Display an error message
      document.getElementById("responseMessage").innerHTML =
        `<p class="text-danger">Erro: ${error.message}</p>`;
    }
  });

// Initialize the map
const map = L.map("map").setView([-15.813, -407.872], 7); // Default to Brasilia coordinates

// Add OpenStreetMap tile layer
L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
  attribution:
    '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
}).addTo(map);

// Marker to display on the map
let marker;

// Variables to store geolocation data
let selectedLatitude = null;
let selectedLongitude = null;
let selectedAddress = "";

// Event listener for map click
map.on("click", async function (e) {
  const lat = e.latlng.lat;
  const lng = e.latlng.lng;

  // Update coordinates display
  document.getElementById("coordinates").textContent =
    `${lat.toFixed(5)}, ${lng.toFixed(5)}`;

  // Reverse geocode the coordinates using Nominatim (OpenStreetMap's geocoding service)
  const response = await fetch(
    `https://nominatim.openstreetmap.org/reverse?lat=${lat}&lon=${lng}&format=json&addressdetails=1`,
  );
  const data = await response.json();

  // Update address display
  const address =
    data && data.display_name ? data.display_name : "No address found";
  document.getElementById("address").textContent = address;

  // Store geolocation data
  selectedLatitude = lat;
  selectedLongitude = lng;
  selectedAddress = address;

  // Fill in the address fields on the form
  fillAddressFields(data);

  // Remove previous marker (if any) and add a new one at the clicked location
  if (marker) {
    marker.remove();
  }

  marker = L.marker([lat, lng])
    .addTo(map)
    .bindPopup(`Address: ${address}`)
    .openPopup();
});

// Function to fill the address fields in the form
function fillAddressFields(data) {
  const address = data.address || {};

  // Fill the address fields in the form
  document.getElementById("logradouro").value = address.road || "";
  document.getElementById("bairro").value = address.neighborhood || "";
  document.getElementById("cidade").value =
    address.city || address.town || address.village || "";
  document.getElementById("estado").value = address.state || "";
  document.getElementById("pais").value = address.country || "";

  // Handle the case if there are no address details available
  if (!address.road) {
    document.getElementById("logradouro").value = "";
  }
}

// Handle the "Send Data" button click
document
  .getElementById("sendDataBtn")
  .addEventListener("click", async function () {
    // Check if coordinates and address are selected
    if (selectedLatitude === null || selectedLongitude === null) {
      document.getElementById("responseMessage").innerHTML =
        `<p class="text-danger">Please click on the map to select a location first.</p>`;
      return;
    }

    // Fill form fields with selected data
    document.getElementById("geolocalizacao").value =
      `${selectedLatitude}, ${selectedLongitude}`;
    document.getElementById("logradouro").value = selectedAddress;

    // Optionally, show success message
    document.getElementById("responseMessage").innerHTML =
      `<p class="text-success">Form populated with map data!</p>`;
  });
