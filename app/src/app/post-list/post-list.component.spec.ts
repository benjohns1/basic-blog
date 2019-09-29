import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PageNotFoundComponent } from '../page-not-found/page-not-found.component';
import { PostComponent } from '../post/post.component';
import { PostListComponent } from './post-list.component';
import { LoginFormComponent } from '../login-form/login-form.component';
import { MatListModule } from '@angular/material/list';
import { AppRoutingModule } from '../app-routing.module';

describe('PostListComponent', () => {
  let component: PostListComponent;
  let fixture: ComponentFixture<PostListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        PostListComponent,
        LoginFormComponent,
        PostComponent,
        PageNotFoundComponent
      ],
      imports: [AppRoutingModule, MatListModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PostListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
