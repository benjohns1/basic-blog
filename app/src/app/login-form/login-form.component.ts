import { Component, OnInit } from '@angular/core';
import { Router } from "@angular/router"
import { AuthenticationService } from '../authentication.service';

@Component({
  selector: 'app-login',
  templateUrl: './login-form.component.html',
  styleUrls: ['./login-form.component.less']
})
export class LoginFormComponent {

  public username = "";
  public password = "";

  constructor(private authenticationService: AuthenticationService, private router: Router) { }

  public login() {
    this.authenticationService.login(this.username, this.password).subscribe(response => {
      if (response === true) {
        this.router.navigate(['/posts'])
      } else {
        // @TODO: display form error message
        console.error(response)
      }
    }, err => {
        console.error(err);
    });
  }

}
