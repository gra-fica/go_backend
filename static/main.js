
async function main(){
	let form_name = document.getElementById("new-name");
	if (form_name == null){
		console.error("something is wrong");
	}
	form_name.focus()
}

main()
