<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			token: null,
			myName: null,
			mySurname:null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.token = localStorage.getItem("Token")
				
				let response = await this.$axios.get("/profile",{headers:{"Token":this.token}});
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async change(){
			this.loading = true;
			this.errormsg = null;
			try {
				this.token = localStorage.getItem("Token")
				let postData = {
					name:this.myName,
					surname:this.mySurname,
				}

				let response = await this.$axios.post("/profile",postData,{headers:{"Token":this.token}});
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>

			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<div v-if="this.some_data" style="display: flex; gap: 40px;">
			<h2>Username</h2>
			<h2>Id</h2>
		</div>
		<div v-if="this.some_data" style="display: flex; gap: 40px;">
			<h2>{{this.some_data.Username}}</h2>
			<h2>{{this.some_data.Id}}</h2>
		</div>
		<div v-if="this.some_data" style="display: flex; gap: 40px;">
			<h2>Name</h2>
			<h2>Surname</h2>
		</div>
		<div v-if="this.some_data" style="display: flex; gap: 20%;">
			<h2>{{this.some_data.Name}}</h2>
			<h2>{{this.some_data.Surname}}</h2>
		</div>
		<div v-if="this.some_data" style="display: flex; gap: 20%;">
			<input type="text" placeholder="Name" v-model="myName" >
			<input type="text" placeholder="Surname" v-model="mySurname" >
			<input type="submit" value="Change" @click="change">
		</div>
		<div v-if="this.some_data" class="grid-container">
			<div v-for="item in this.some_data.Images" >
				<div class="grid-item">
					<img v-bind:src="'http://localhost:3000'+item.Location">
					<div style="display: flex;gap: 20%;">
						<h5>{{item.Text}}</h5>
						<h4>{{item.Comments}}</h4>
						<h4>{{item.Likes}}</h4>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.grid-container {
  display: grid;
  grid-template-columns: auto auto auto;
  padding: 10px;
}
.grid-item {
  padding: 20px;
  font-size: 30px;
  text-align: center;
  border: 1px solid #000000;
}
.grid-item img{
	object-fit: cover;
	width: 230px;
  	height: 230px;
}
</style>