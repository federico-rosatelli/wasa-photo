<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			token_data: "",
			myUsername: ""
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			//this.myUsername = "";
			let dataPost = {
				username:this.myUsername
			}
			const headers = {
				'Content-Type': 'application/json'
			};
			// this.$axios.post("/signin",{
			// 	username:this.myUsername
			// 	})
			// 	.then((response)=>{
			// 		console.log(response);
			// 		this.token_data = response.data;
			// 	})
			// 	.catch((e)=>{
			// 		this.errormsg = e.toString();
			// 	})
			let response = await this.$axios.post("/signin",dataPost)
			this.token_data = response.data
			// fetch("http://127.0.0.1:3000/signin",{
			// 	method:'POST',
			// 	headers:{
			// 		'Content-Type':'application/json'
			// 	},
			// 	body: JSON.stringify(dataPost)
			// 	})
			// 	.then((response) => response.text())
  			// 	.then((data) => console.log(data));
				// .catch((error)=>{
				// 	this.errormsg = error.toString();
				// })
			// try {
            //     //let username = myUsername
			// 	let response = await this.$axios.post("/signin",dataPost);
			// 	//let response = await this.$axios.get("/")
			// 	this.token_data = response.data;
			// } catch (e) {
			// 	this.errormsg = "prova "+e.toString();
			// }
			this.loading = false;
		},
	},
	// mounted() {
	// 	this.refresh()
	// }
}
</script>

<template>
	<div>
		<div class="login-form">
			<h1>Login Form</h1>
			    <input type="text" placeholder="Username" v-model="myUsername" >
				<input type="submit" value="Log In" @click="refresh">
			<div style="padding:15px">
		</div>
    </div>
	<h1>{{this.token_data}}</h1>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
