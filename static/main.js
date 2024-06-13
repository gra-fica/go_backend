
async function main(){
	/** 
		* @type HTMLFormElement
	*/
	let form = document.getElementById("new-item-form");
	if(form == null){
		console.log("form not found!!");
		return;	
	}

	form.addEventListener("keydown", ev => {
		if(ev.key == "Enter"){
			htmx.ajax('POST', '/api/v1/htmx/search/product', {
				target: "#search-results",
			})
		}
	});
}

main()
