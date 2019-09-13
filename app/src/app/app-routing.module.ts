import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { PostListComponent } from 'src/app/post-list/post-list.component';
import { PageNotFoundComponent } from 'src/app/page-not-found/page-not-found.component';
import { LoginFormComponent } from './login-form/login-form.component';
import { BlogPostFilter } from './blog.service';
import { PostComponent } from './post/post.component';


const routes: Routes = [
  { path: 'login', component: LoginFormComponent },
  { path: 'posts', component: PostListComponent },
  { path: 'posts/deleted', component: PostListComponent, data: { filter: BlogPostFilter.Deleted } },
  { path : 'post', redirectTo: '/post/new', pathMatch: 'full' },
  { path : 'post/new', component: PostComponent, data: { new: true, edit: true } },
  { path : 'post/:id', component: PostComponent },
  { path : 'post/:id/edit', component: PostComponent, data: { edit: true } },
  { path: '', redirectTo: '/posts', pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
