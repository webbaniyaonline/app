$('.hrefModal').click(function(){
    //alert("Modal Box hrefModal")
    var modalUrls=$(this).attr('data-href');
    var tid=$(this).attr('data-tid');

    //alert(modalUrls)

    $('#transModal').modal('show');
    $('#transModal').modal('show').find('.modal-body').load(modalUrls);
    $('#transModal .modal-dialog').css({"max-width":"80%", "margin-top": "20px"});
    $('#transModal .modal-title').html("Transaction Details - " + tid);
});

$(".loaderx").click(function(){
    //alert(1111)
$(".loaderx").html("<i class='fa-solid fa-spinner fa-spin-pulse'></i>");
});


// For Checkout Page JS

// Function to format time as MM:SS
function formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

// The jQuery function from Step 1 goes here
function getUrlParameter(name) {
            name = name.replace(/[\[]/, '\\[').replace(/[\]]/, '\\]');
            var regex = new RegExp('[\\?&]' + name + '=([^&#]*)');
            var results = regex.exec(window.location.href);
            return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '));
        }

        $('#replybtn').click(function(){

          location.reload(true);
      });
    
      $("#crypto_network").change(function(){
        
        var myTid=$("#crypto_network option:selected").attr("data-id");
        $('#crypto_id').val(myTid);
    });


      function fetchStates() {
        //alert(111111111)
        const crypto_code = document.getElementById("crypto_code").value;
        const stateSelect = document.getElementById("crypto_network");

   // Get the dropdown element by its ID
    var cryptoDropdown = document.getElementById("crypto_code");
    // Get the selected option
    var selectedOption = cryptoDropdown.options[cryptoDropdown.selectedIndex];
    // Get the value of the custom attribute
    var dataCode = selectedOption.getAttribute("data-title");
    document.getElementById("crypto_title").value = dataCode;

        // Clear existing options
        stateSelect.innerHTML = '<option value="" >Select Network </option>';
    
        if (crypto_code) {
            fetch(`/get-network?crypto_code=${crypto_code}`)
            .then(response => response.json())
            .then(networks => {
                networks.forEach(network => {
                    const option = document.createElement("option");
                    option.value = network.Crypto_network_short;
                    option.setAttribute('data-id', network.Crypto_id)
                    option.textContent = network.Crypto_network_short;
                    stateSelect.appendChild(option);
                    
                });
            })
                .catch(error => console.error('Error fetching states:', error));
        }
    }


    function fetchnetwork() {
        //alert(111111111)
        const crypto_code = document.getElementById("crypto_code").value;
        const stateSelect = document.getElementById("crypto_network");

   // Get the dropdown element by its ID
    var cryptoDropdown = document.getElementById("crypto_code");
    // Get the selected option
    var selectedOption = cryptoDropdown.options[cryptoDropdown.selectedIndex];
    // Get the value of the custom attribute
    var dataCode = selectedOption.getAttribute("data-title");
    document.getElementById("crypto_title").value = dataCode;

        // Clear existing options
        stateSelect.innerHTML = '<option value="" >Select Network </option>';
    
        if (crypto_code) {
            fetch(`/get-network?crypto_code=${crypto_code}&tab=coin_list`)
            .then(response => response.json())
            .then(networks => {
                networks.forEach(network => {
                    const option = document.createElement("option");
                    option.value = network.Coin_network;
                    option.setAttribute('data-id', network.Coin_id)
                    option.textContent = network.Coin_network;
                    stateSelect.appendChild(option);
                    
                });
            })
                .catch(error => console.error('Error fetching states:', error));
        }
    }

//For change Date Format
    // Function to format the date
function formatDate(dateStr) {
    const date = new Date(dateStr);
    const options = {
        year: 'numeric',
        month: 'numeric', //long
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',  // Include seconds if needed
        hour12: true         // Set to false for 24-hour format
    };
    return date.toLocaleDateString(undefined, options);
}

// Get all date cells
const dateCells = document.querySelectorAll('.date-cell');

// Loop through each date cell and update its content
dateCells.forEach(cell => {
    const originalDate = cell.textContent;
    const formattedDate = formatDate(originalDate);
    cell.textContent = formattedDate;
});

//End change Date Format

// copy Code
function CopyToClipboard(text) {
    text=$(text).text();
    
    var $txt = $('<textarea />');

    $txt.val(text)
        .css({ width: "1px", height: "1px" })
        .appendTo('body');

    $txt.select();
    

    if (document.execCommand('copy')) {
        $txt.remove();
        alert("Copied");
    }
}
//Get currency sumbol from currency
function getCurrencySymbol(currencyCode) {
    alert(33333)
    // Create a new instance of Intl.NumberFormat
    const formatter = new Intl.NumberFormat('en', {
      style: 'currency',
      currency: currencyCode,
      minimumFractionDigits: 0, // Remove decimal places for the symbol
    });
  
    // Format a value and extract the currency symbol
    const parts = formatter.formatToParts(0);
    const symbol = parts.find(part => part.type === 'currency').value;

    alert(symbol)
  
    return symbol;
  }

      //For Status by id
    // Function to status Data
function statusData(status,tid) {
    var statusVal="";
    var statusHtml="";
    
    if(status==0){
      statusVal="Waiting";
      statusHtml='<span class="bg-warning-600 text-white border-warning-600 px-16 py-2 radius-12 fw-medium text-sm">Waiting</span>';
    }else if(status==1){
      statusVal="FullPay";
      statusHtml='<span class="bg-success-500 text-white border-success-600 px-16 py-2 radius-12 fw-medium text-sm">FullPay</span>';
    }else if(status==2){
      statusVal="OverPay";
      statusHtml='<span class="bg-success-600 text-white border-success-600 px-16 py-2 radius-12 fw-medium text-sm">OverPay</span>';
    }else if(status==3){
      statusVal="UnderPay";
      statusHtml='<span class="bg-success-400 text-white border-success-600 px-16 py-2 radius-12 fw-medium text-sm">UnderPay</span>';
    }else if(status==8){
      statusVal="Decline";
      statusHtml='<span class="bg-danger-600 text-white border-danger-600 px-16 py-2 radius-12 fw-medium text-sm">Decline</span>';
    }else if(status==9){
      statusVal="Dispute";
      statusHtml='<span class="bg-lilac-600 text-white border-lilac-600 px-16 py-2 radius-12 fw-medium text-sm">Dispute</span>';
    }else{
      statusVal="";
      statusHtml="";
    }
    
    if(tid==2){
      return statusHtml;
    }else{
      return statusVal;
    }
      
  }
  
     // Get all date cells
  const paystatus = document.querySelectorAll('.pay-status');
  
  // Loop through each date cell and update its content
  paystatus.forEach(cell => {
      const statusId = cell.textContent;
      const fetchedStatus = statusData(statusId,2);
      cell.innerHTML = fetchedStatus;
  });

// Display Balance on 6 digit decimal number
  // Get all Balance Amount cells
const balanceCells = document.querySelectorAll('.balance-limit');
// Loop through each Balance cell and update its content
balanceCells.forEach(cell => {
    const originalAmount = cell.textContent;
    const formattedAmount = parseFloat(originalAmount).toFixed(6); // set 6 digit decimal number
    cell.textContent = formattedAmount;
});

 // Display Support ticket status
 const supportStatus = document.querySelectorAll('.support-status');
 // Loop through each Balance cell and update its content
 supportStatus.forEach(cell => {
     const originalStatus = cell.textContent;
     var formattedStatus=""
     if(originalStatus==1){
      formattedStatus ="Open";
     }else if(originalStatus==2){
      formattedStatus ="Replied";
     }else{
       formattedStatus ="Closed";
     }
     cell.textContent = formattedStatus;
 });

 // For check used device
 function detectDevice() {
  const userAgent = navigator.userAgent || navigator.vendor || window.opera;

  if (/android/i.test(userAgent)) {
    return "Android";
  } else if (/iPad|iPhone|iPod/.test(userAgent) && !window.MSStream) {
    return "iOS";
  } else {
    return "Web";
  }
}