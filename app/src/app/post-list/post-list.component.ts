import { Component, OnInit } from '@angular/core';
import { BlogService, BlogPostFilter } from '../blog.service';
import { ActivatedRoute } from '@angular/router';
import { Post } from '../post';

@Component({
  selector: 'app-post-list',
  templateUrl: './post-list.component.html',
  styleUrls: ['./post-list.component.less']
})
export class PostListComponent implements OnInit {

  public posts: Post[] = [];
  public filter: BlogPostFilter;

  constructor(private blogService: BlogService, route: ActivatedRoute) {
    route.data.forEach(d => {
      if (this.filter !== d.filter) {
        this.filter = d.filter;
        this.loadPosts();
      }
    });
    this.loadPosts();
  }

  ngOnInit() {
  }

  private loadPosts() {
    this.blogService.getPosts(this.filter).subscribe(posts => {
      this.posts = posts;
    }, error => {
        // @TODO: show friendly error in UI
        console.error(error)
    });
  }

}
