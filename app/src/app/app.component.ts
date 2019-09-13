import { Component } from '@angular/core';
import { AuthenticationService } from './authentication.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less']
})
export class AppComponent {
  title = 'Basic Blog';

  constructor(private authenticationService: AuthenticationService) { }

  public isLoggedIn(): boolean {
    return this.authenticationService.isLoggedIn();
  }

  public logout() {
    this.authenticationService.logout();
  }

}
